package disco

import (
	"Archivos/PY1/estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"
)

type listaEBR struct {
	size   int64
	inicio int64
	antes  int64
	sig    *listaEBR
}

type Fdisk struct {
	Path   string
	Name   string
	Fit    string
	Tipo   string
	Size   int64
	Unit   string
	Add    int64
	Delete string
}

func Administrar(comando Fdisk) {
	mbr := LeerDisco(comando.Path)
	if mbr.Size > 0 {
		if comando.Size > -1 { //crear particion
			crearParticion(comando, &mbr)
		} else if comando.Add != 0 {
			modificarSize(comando, &mbr)
		} else if comando.Delete != "" {
			fmt.Println("Esta seguro que desea borrar la particion?")
			fmt.Println("1 Si              2 No")
			op := ""
			fmt.Scanln(&op)
			if op == "1" {
				borrada := borrarPart(comando, &mbr)
				if borrada.Status != GetChar("0") {
					var u Unmount
					u.Desmontadas = append(u.Desmontadas, comando.Name)
					if borrada.Tipo == GetChar("p") || borrada.Tipo == GetChar("e") {
						Desmontar(u)
						if strings.EqualFold(comando.Delete, "full") {
							f, err := os.OpenFile(comando.Path, os.O_RDWR, 0755) //leer o escribir
							if err != nil {
								log.Fatal(err)
							}
							f.Seek(int64(borrada.Start), 0)
							var p estructuras.Particion
							var b bytes.Buffer
							binary.Write(&b, binary.BigEndian, &p)
							EscribirBytes(f, b.Bytes())
							f.Close()
							fmt.Println("Particion borrada con exito")
						}
					} else if borrada.Tipo == GetChar("l") {
						Desmontar(u)
						fmt.Println("Particion borrada con exito")
					}
				} else {
					fmt.Println("Error: La particion no existe")
				}
			}
		}
		file, err := os.OpenFile(comando.Path, os.O_RDWR, 0755) //leer o escribir
		if err != nil {
			log.Fatal(err)
		}
		file.Seek(0, 0)
		binario := new(bytes.Buffer)
		binary.Write(binario, binary.BigEndian, &mbr)
		EscribirBytes(file, binario.Bytes())
		file.Close()
	}
}

func crearParticion(comando Fdisk, mbr *estructuras.MBR) {
	if int64(comando.Size) < (mbr.Size - int64(unsafe.Sizeof(*mbr))) {
		full := ComprobarParticionesLlenas(mbr)
		if !comprobarName(mbr, comando.Name) {
			if ComprobarParticionesVacias(mbr) {
				if strings.EqualFold(comando.Tipo, "p") { //primaria
					mbr.Particiones[0] = nuevaParticion(comando)
				} else if strings.EqualFold(comando.Tipo, "e") { //extendida
					mbr.Particiones[0] = nuevaParticion(comando)
					ebr := nuevoEBR(comando)
					ebr.Start = mbr.Particiones[0].Start
					escribirEBR(comando, &ebr)
				} else { //logica
					fmt.Println("Error: No existe particion extendida donde pueda agregar una logica")
					return
				}
			} else if full { // 4 particiones llenas
				if strings.EqualFold(comando.Tipo, "l") {
					extendida := GetExtendida(mbr)
					if extendida.Size > 0 {
						adminEBR(extendida, comando)
					} else {
						fmt.Println("Error: no se puede crear una particion logica, no existe una extendida")
						return
					}
				} else {
					fmt.Println("Error: No se puede crear la particion, ha llegado al numero maximo de particiones")
					return
				}
			} else {
				switch comando.Tipo {
				case "p":
					fFPart(mbr, comando)
				case "e":
					extendida := GetExtendida(mbr)
					if extendida.Tipo != GetChar("e") {
						fFPart(mbr, comando)
						extendida = GetExtendida(mbr)
						ebr := nuevoEBR(comando)
						ebr.Start = extendida.Start
						escribirEBR(comando, &ebr)
					} else {
						fmt.Println("Error: Ya existe una particion extendida, no puede crear otra")
						return
					}
				case "l":
					extendida := GetExtendida(mbr)
					if extendida.Tipo == GetChar("e") {
						adminEBR(extendida, comando)
					} else {
						fmt.Println("Error: No puede crear una particion logica sin haber una extendida")
						return
					}
				}
			}
		} else {
			fmt.Println("Error: La particion ya existe")
		}

	} else {
		fmt.Println("Error: No hay suficiente espacio en el disco para crear la particion")
	}
}

func modificarSize(comando Fdisk, mbr *estructuras.MBR) {
	var name [16]byte
	copy(name[:], comando.Name)
	if comando.Add > 0 { //sumar espacio a la particion
		var part estructuras.Particion
		var aux estructuras.Particion
		i := 0
		for i = 0; i < 4; i++ {
			if mbr.Particiones[i].Tipo == GetChar("e") {
				aux = mbr.Particiones[i]
			}
			if name == mbr.Particiones[i].Name {
				part = mbr.Particiones[i]
				break
			}
		}
		if i == 3 {
			disponible := mbr.Size - (part.Start + part.Size)
			if disponible >= comando.Add {
				part.Size += comando.Add
				mbr.Particiones[i] = part
			} else {
				fmt.Println("Error: No hay suficiente espacio para agrandar la particion")
			}
		} else if i < 3 {
			disponible := mbr.Particiones[i+1].Start - (part.Start + part.Size)
			if disponible > comando.Add {
				part.Size += comando.Add
				mbr.Particiones[i] = part
			} else {
				fmt.Println("Error: No hay suficiente espacio para agrandar la particion")
			}
		} else if i == 4 {
			if aux.Tipo == GetChar("e") { // verificar si hay extendida
				var ebr estructuras.EBR
				ebr = LeerEBR(comando.Path, aux.Start)
				for ebr.Next != -1 && ebr.Name != name {
					ebr = LeerEBR(comando.Path, ebr.Next)
				}
				if ebr.Name == name {
					if ebr.Next == -1 {
						disponible := aux.Size - (ebr.Size + ebr.Start)
						if disponible > comando.Add {
							ebr.Size += comando.Add
							escribirEBR(comando, &ebr)
						} else {
							fmt.Println("Error: No hay espacio suficiente para agregarle a la particion")
						}
					} else {
						disponible := ebr.Next - (ebr.Size + ebr.Start)
						if disponible > comando.Add {
							ebr.Size += comando.Add
							escribirEBR(comando, &ebr)
						} else {
							fmt.Println("Error: No hay espacio suficiente para agregarle a la particion")
						}
					}
				} else {
					fmt.Println("Error: No existe la particion")
				}
			} else {
				fmt.Println("Error: No existe la particion")
			}
		}

	} else { //quitarle espacio a la particion
		var part estructuras.Particion
		for i := 0; i < 4; i++ {
			if name == mbr.Particiones[i].Name {
				if (-1 * comando.Add) < mbr.Particiones[i].Size {
					mbr.Particiones[i].Size += comando.Add

				} else {
					fmt.Println("Error: El espacio a eliminar es mayor al que ocupa la particion")
				}
				return
			}
			if mbr.Particiones[i].Tipo == GetChar("e") {
				part = mbr.Particiones[i]
			}
		}
		if part.Tipo == GetChar("e") {
			var ebr estructuras.EBR
			ebr = LeerEBR(comando.Path, part.Start)
			for ebr.Next != 1 && ebr.Name != name {
				ebr = LeerEBR(comando.Path, ebr.Next)
			}
			if ebr.Name == name {
				if (-1 * comando.Add) < ebr.Size {
					ebr.Size += comando.Add
					escribirEBR(comando, &ebr)
				} else {
					fmt.Println("Error: El espacio a eliminar es mayor al que ocupa la particion")
				}
			} else {
				fmt.Println("Error: La particion no existe")
			}
		} else {
			fmt.Println("Error: La particion no existe")
		}
	}
}

func borrarPart(c Fdisk, mbr *estructuras.MBR) estructuras.Particion {
	var name [16]byte
	copy(name[:], c.Name)
	var part estructuras.Particion
	var aux estructuras.Particion
	for i := 0; i < 4; i++ {
		if name == mbr.Particiones[i].Name {
			part = mbr.Particiones[i]
			mbr.Particiones[i].Size = 0
			mbr.Particiones[i].Start = mbr.Size
			mbr.Particiones[i].Status = GetChar("0")
			mbr.Particiones = ordenarParticiones(mbr.Particiones, 4)
			if part.Tipo == GetChar("e") {
				var ebr estructuras.EBR
				ebr = LeerEBR(c.Path, part.Start)
				ebr = nuevoEBR(c)
				escribirEBR(c, &ebr)
			}
			return part
		}
		if mbr.Particiones[i].Tipo == GetChar("e") {
			aux = mbr.Particiones[i]
		}
	}
	if aux.Tipo == GetChar("e") {
		aux.Tipo = GetChar("l")
		var ebr estructuras.EBR
		var ebraux estructuras.EBR
		ebr = LeerEBR(c.Path, aux.Start)
		ebraux = ebr
		if name == ebr.Name {
			ebr.Size = 0
			ebr.Status = 0
			for i := 0; i < 16; i++ {
				ebr.Name[i] = GetChar("0")
			}
			escribirEBR(c, &ebr)
			aux.Name = ebraux.Name
			return aux
		}
		for ebr.Next != -1 && ebr.Name != name {
			ebraux = ebr
			ebr = LeerEBR(c.Path, ebr.Next)
		}
		if name == ebr.Name {
			ebraux.Next = ebr.Next
			var escrito estructuras.EBR
			escrito.Next = 0
			f, err := os.OpenFile(c.Path, os.O_RDWR, 0755) //leer o escribir
			if err != nil {
				log.Fatal(err)
			}
			f.Seek(int64(ebr.Start), 0)
			var ebrb bytes.Buffer
			binary.Write(&ebrb, binary.BigEndian, &escrito)
			EscribirBytes(f, ebrb.Bytes())
			f.Close()
			aux.Name = ebr.Name
			return aux
		}
		aux.Status = GetChar("0")
		return aux
	}
	return aux
}

func LeerDisco(path string) estructuras.MBR {
	mbr := estructuras.MBR{}
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error: El disco no existe, no se puede crear la particion")
		return mbr
	}
	defer file.Close()
	var size int = int(unsafe.Sizeof(mbr))
	m := LeerBytes(file, size)
	buffer := bytes.NewBuffer(m)
	err = binary.Read(buffer, binary.BigEndian, &mbr)
	if err != nil {
		log.Fatal(err)
	}
	return mbr
}

func LeerEBR(path string, seek int64) estructuras.EBR {
	var ebr estructuras.EBR
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error: El disco no existe")
		return ebr
	}
	file.Seek(seek, 0)
	e := LeerBytes(file, int(unsafe.Sizeof(ebr)))
	buffer := bytes.NewBuffer(e)
	err = binary.Read(buffer, binary.BigEndian, &ebr)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	return ebr
}

func adminEBR(extendida estructuras.Particion, comando Fdisk) {
	ebr := LeerEBR(comando.Path, extendida.Start)
	if ebr.Next == -1 && ebr.Status == GetChar("0") {
		ebr.Status = GetChar("1")
		ebr.Fit = GetChar(comando.Fit)
		ebr.Start = extendida.Start
		ebr.Size = comando.Size
		ebr.Next = -1
		copy(ebr.Name[:], comando.Name)
		f, err := os.OpenFile(comando.Path, os.O_RDWR, 0755) //leer o escribir
		if err != nil {
			fmt.Println(err)
			return
		}
		f.Seek(int64(extendida.Start), 0)
		var ebrb bytes.Buffer
		binary.Write(&ebrb, binary.BigEndian, &ebr)
		EscribirBytes(f, ebrb.Bytes())
		f.Close()
	} else {
		//agregarEBR(extendida, comando, ebr)
	}
}

func fFPart(mbr *estructuras.MBR, c Fdisk) {
	var ini [4]int64
	var disponible [4]int64
	posicion := 0
	part := nuevaParticion(c)
	part.Status = GetChar("0")
	tam := mbr.Particiones[0].Start - int64(unsafe.Sizeof(*mbr))
	if tam > 0 {
		posicion++
		ini[0] = int64(unsafe.Sizeof(*mbr))
		disponible[0] = tam
	}
	tam = 0
	for i := 0; i < 3; i++ {
		if mbr.Particiones[i+1].Status != GetChar("1") || mbr.Particiones[i+1].Start == mbr.Size {
			ini[posicion] = mbr.Particiones[i].Start + mbr.Particiones[i].Size
			disponible[posicion] = mbr.Size - ini[posicion]
			posicion++
			break
		}
		ini[posicion] = mbr.Particiones[i].Start + mbr.Particiones[i].Size
		disponible[posicion] = mbr.Particiones[i+1].Start - ini[posicion]
		posicion++
	}
	for i := 0; i < posicion; i++ {
		if disponible[i] >= c.Size {
			part.Start = ini[i]
			part.Status = GetChar("1")
			mbr.Particiones[posicion] = part
			break
		}
	}
	if part.Status != GetChar("1") {
		fmt.Println("Error: No hay espacio suficiente para crear la particion")
	} else {
		mbr.Particiones = ordenarParticiones(mbr.Particiones, posicion)
	}
}

/*
func peorAjuste(ebr, ebraux estructuras.EBR, p, u *listaEBR, c Fdisk) {
	if u.size >= c.Size {
		if u.inicio == u.antes {
			ebr = LeerEBR(c.Path, u.inicio)
			ebr.Status = GetChar("1")
			ebr.Fit = GetChar(c.Fit)
			ebr.Start = u.inicio
			ebr.Size = c.Size
			for i := 0; i < 16; i++ {
				ebr.Name[i] = c.Name[i]
			}
			escribirEBR(c, &ebr)
		} else {
			ebr = LeerEBR(c.Path, u.antes)
			auxini := ebr.Next
			ebraux.Status = GetChar("1")
			ebraux.Start = u.inicio
			ebraux.Size = c.Size
			ebraux.Next = auxini
			for i := 0; i < 16; i++ {
				ebr.Name[i] = c.Name[i]
			}
			ebr.Next = ebraux.Start
			escribirEBR(c, &ebr)
			escribirEBR(c, &ebraux)
		}
	} else {
		fmt.Println("Error: No hay espacio suficiente para crear la particion")
	}
}*/

func ajustar(ebr, ebraux estructuras.EBR, p, u *listaEBR, c Fdisk) {
	var aux *listaEBR
	aux = p
	for aux != u.sig {
		if aux.size >= c.Size {
			if aux.antes == aux.inicio {
				ebr = LeerEBR(c.Path, aux.inicio)
				ebr.Status = GetChar("1")
				ebr.Fit = GetChar(c.Fit)
				ebr.Start = aux.inicio
				ebr.Size = c.Size
				for i := 0; i < 16; i++ {
					ebr.Name[i] = c.Name[i]
				}
				escribirEBR(c, &ebr)
			} else {
				ebr = LeerEBR(c.Path, aux.antes)
				auxini := ebr.Next
				ebraux.Status = GetChar("1")
				ebraux.Fit = GetChar(c.Fit)
				ebraux.Start = aux.inicio
				ebraux.Size = c.Size
				ebraux.Next = auxini
				for i := 0; i < 16; i++ {
					ebr.Name[i] = c.Name[i]
				}
				ebr.Next = ebraux.Start
				escribirEBR(c, &ebr)
				escribirEBR(c, &ebraux)
			}
		}
		aux = aux.sig
	}
	if aux == u.sig {
		fmt.Println("Error: No hay espacio suficiente para crear la particion")
	}
}

func ordenarParticiones(parts [4]estructuras.Particion, pos int) [4]estructuras.Particion {
	var aux estructuras.Particion
	for i := 1; i < pos; i++ {
		for j := 0; j < pos; j++ {
			if parts[j].Start > parts[j+1].Start {
				aux = parts[j]
				parts[j] = parts[j+1]
				parts[j+1] = aux
			}
		}
	}
	return parts
}

func agregarEBR(exte estructuras.Particion, c Fdisk, ebr estructuras.EBR) {
	var primero listaEBR
	var ultimo listaEBR
	var aux estructuras.EBR
	var size int64
	var pos int64
	if ebr.Status == GetChar("0") && ebr.Next != -1 {
		aux = LeerEBR(c.Path, ebr.Next)
		size = aux.Start - ebr.Start
		pos = ebr.Start
		listarEBR(&primero, &ultimo, size, pos, pos)
		ebr = aux
	}
	for ebr.Next != -1 {
		aux = LeerEBR(c.Path, ebr.Next)
		pos = ebr.Size + ebr.Start
		size = aux.Start - pos
		listarEBR(&primero, &ultimo, size, pos, ebr.Start)
		ebr = aux
	}
	if ebr.Next == -1 {
		size = exte.Size + exte.Start - (ebr.Start + ebr.Size)
		pos = ebr.Start + ebr.Size
		listarEBR(&primero, &ultimo, size, pos, ebr.Start)
	}
	switch exte.Fit {
	case GetChar("f"):
		ajustar(ebr, aux, &primero, &ultimo, c)
	case GetChar("b"):
		ordenarEBR(&primero, &ultimo)
	default:
		ordenarEBR(&primero, &ultimo)
		//peorAjuste(ebr, aux, &primero, &ultimo, c)
	}
}

func listarEBR(p, u *listaEBR, size, actual, anterior int64) {
	if size > 0 {
		var nuevo *listaEBR
		var nuevoaux listaEBR
		nuevoaux.size = size
		nuevoaux.antes = anterior
		nuevoaux.inicio = actual
		nuevo = &nuevoaux
		if p == nil {
			p = nuevo
			u = nuevo
		} else {
			u.sig = nuevo
			u = nuevo
		}
	}
}

func ordenarEBR(p, u *listaEBR) { //metodo que ordena para encontrar el mejor ajuste
	if p.inicio != 0 {
		var aux *listaEBR
		var actual *listaEBR
		auxsize := int64(0)
		auxini := int64(0)
		auxantes := int64(0)
		aux = p
		for aux != u {
			actual = aux.sig
			for actual.inicio != 0 {
				if aux.size > actual.size {
					auxsize = aux.size
					auxini = aux.inicio
					auxantes = aux.antes
					aux.size = actual.size
					aux.inicio = actual.inicio
					aux.antes = actual.antes
					actual.antes = auxantes
					actual.inicio = auxini
					actual.size = auxsize
				}
				actual = actual.sig
			}
			aux = aux.sig
		}
	}
}

func escribirEBR(comando Fdisk, ebr *estructuras.EBR) {
	f, err := os.OpenFile(comando.Path, os.O_RDWR, 0755) //leer o escribir
	if err != nil {
		fmt.Println(err)
		return
	}
	f.Seek(int64(ebr.Start), 0)
	var ebrb bytes.Buffer
	binary.Write(&ebrb, binary.BigEndian, ebr)
	EscribirBytes(f, ebrb.Bytes())
	f.Close()
}

func GetExtendida(mbr *estructuras.MBR) estructuras.Particion {
	var part estructuras.Particion
	for i := 0; i < 4; i++ {
		if mbr.Particiones[i].Tipo == GetChar("e") {
			return mbr.Particiones[i]
		}
	}
	return part
}

func ComprobarParticionesVacias(mbr *estructuras.MBR) bool {
	for i := 0; i < 4; i++ {
		if mbr.Particiones[i].Status == GetChar("1") {
			return false
		}
	}
	return true
}

func ComprobarParticionesLlenas(mbr *estructuras.MBR) bool {
	for i := 0; i < 4; i++ {
		if mbr.Particiones[i].Status != GetChar("1") {
			return false
		}
	}
	return true
}

func comprobarName(mbr *estructuras.MBR, name string) bool {
	for i := 0; i < 4; i++ {
		var n [16]byte
		copy(n[:], name)
		if n == mbr.Particiones[i].Name {
			return true
		}
	}
	return false
}

func nuevaParticion(c Fdisk) estructuras.Particion {
	var mbr estructuras.MBR
	var part estructuras.Particion
	part.Size = c.Size
	copy(part.Name[:], c.Name)
	part.Status = GetChar("1")
	part.Fit = GetChar(c.Fit)
	part.Tipo = GetChar(c.Tipo)
	part.Start = int64(unsafe.Sizeof(mbr))
	return part
}

func nuevoEBR(c Fdisk) estructuras.EBR {
	var ebr estructuras.EBR
	ebr.Fit = GetChar("w")
	ebr.Next = -1
	ebr.Size = 0
	ebr.Status = GetChar("0")
	return ebr
}

func GetChar(s string) byte {
	return s[0]
}

func LeerBytes(f *os.File, n int) []byte {
	bytes := make([]byte, n)
	_, err := f.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
