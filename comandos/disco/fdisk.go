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
	size   int
	inicio int
	antes  int
	sig    *listaEBR
}

type Fdisk struct {
	Path   string
	Name   string
	Fit    string
	Tipo   string
	Size   int
	Unit   string
	Add    int
	Delete string
}

func administrar(comando Fdisk) {
	mbr := leerDisco(comando.Path)
	if mbr.Size > 0 {
		if comando.Size > -1 { //crear particion
			if comando.Size > (mbr.Size - int(unsafe.Sizeof(mbr))) {
				full := comprobarParticionesLlenas(mbr)
				if comprobarName(mbr, comando.Name) {
					if comprobarParticionesVacias(mbr) {
						if strings.EqualFold(comando.Tipo, "p") { //primaria
							mbr.Particiones[0] = nuevaParticion(comando, mbr)
						} else if strings.EqualFold(comando.Tipo, "e") { //extendida
							mbr.Particiones[0] = nuevaParticion(comando, mbr)
							ebr := nuevoEBR(comando)
							ebr.Start = mbr.Particiones[0].Start
							escribirEBR(comando, &ebr)
						} else { //logica
							fmt.Println("Error: No existe particion extendida donde pueda agregar una logica")
							return
						}
					} else if full { // 4 particiones llenas
						if strings.EqualFold(comando.Tipo, "l") {
							extendida := getExtendida(mbr)
							if extendida.Size > 0 {
								adminEBR(extendida, comando)
							} else {
								fmt.Println("Error: no se puede crear una particion logica, no existe una extendida")
							}
						} else {
							fmt.Println("Error: No se puede crear la particion, ha llegado al numero maximo de particiones")
						}
					} else {
						switch comando.Tipo {
						case "p":
							mbr = fFPart(mbr, comando)
						case "e":
							extendida := getExtendida(mbr)
							if extendida.Tipo != getChar("e") {
								mbr = fFPart(mbr, comando)
								extendida = getExtendida(mbr)
								ebr := nuevoEBR(comando)
								ebr.Start = extendida.Start
								escribirEBR(comando, &ebr)
							} else {
								fmt.Println("Error: Ya existe una particion extendida, no puede crear otra")
							}
						case "l":
							extendida := getExtendida(mbr)
							if extendida.Tipo == getChar("e") {
								adminEBR(extendida, comando)
							} else {
								fmt.Println("Error: No puede crear una particion logica sin haber una extendida")
							}
						}
					}

				}

			} else {
				fmt.Println("Error: No hay suficiente espacio en el disco para crear la particion")
			}

		}
	}
}

func crearParticion(comando Fdisk) {
	mbr := leerDisco(comando.Path)
	size := 0
	if mbr.Size > 0 {
		fi, err := os.Stat(comando.Path)
		if err != nil {
			log.Fatal(err)
		}
		size = int(fi.Size())
		if size > comando.Size {

		} else {
			fmt.Println("Error: espacio insuficiente en el disco para crear la particion")
		}
	}
}

func modificarSize(comando Fdisk) {

}

func borrarPart(comando Fdisk) {

}

func escribirEnDisco() {

}

func leerDisco(path string) estructuras.MBR {
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

func leerEBR(f *os.File) estructuras.EBR {
	var ebr estructuras.EBR
	e := LeerBytes(f, int(unsafe.Sizeof(ebr)))
	buffer := bytes.NewBuffer(e)
	err := binary.Read(buffer, binary.BigEndian, &ebr)
	if err != nil {
		log.Fatal(err)
	}
	return ebr
}

func adminEBR(extendida estructuras.Particion, comando Fdisk) {
	f, err := os.OpenFile(comando.Path, os.O_RDWR, 0755) //leer o escribir
	if err != nil {
		fmt.Println(err)
		return
	}
	f.Seek(int64(extendida.Start), 0)
	ebr := leerEBR(f)
	if ebr.Next == -1 && ebr.Status == getChar("0") {
		ebr.Status = getChar("1")
		ebr.Fit = getChar(comando.Fit)
		ebr.Start = extendida.Start
		ebr.Size = comando.Size
		ebr.Next = -1
		for i := 0; i < 16; i++ {
			ebr.Name[i] = comando.Name[i]
		}
		f.Seek(int64(extendida.Start), 0)
		var ebrb bytes.Buffer
		binary.Write(&ebrb, binary.BigEndian, &ebr)
		EscribirBytes(f, ebrb.Bytes())
		f.Close()
	} else {
		agregarEBR(f, extendida, comando, ebr)
	}
}

func fFPart(mbr estructuras.MBR, c Fdisk) estructuras.MBR {
	var ini [4]int
	var disponible [4]int
	posicion := 0
	part := nuevaParticion(c, mbr)
	part.Status = getChar("0")
	tam := mbr.Particiones[0].Start - int(unsafe.Sizeof(mbr))
	if tam > 0 {
		posicion++
		ini[0] = int(unsafe.Sizeof(mbr))
		disponible[0] = tam
	}
	tam = 0
	for i := 0; i < 3; i++ {
		if mbr.Particiones[i+1].Status == getChar("0") || mbr.Particiones[i+1].Start == mbr.Size {
			ini[posicion] = mbr.Particiones[i].Start + mbr.Particiones[i].Size
			disponible[posicion] = mbr.Size - ini[posicion]
			posicion++
			break
		}
		ini[posicion] = mbr.Particiones[i].Start + mbr.Particiones[i].Size
		disponible[posicion] = mbr.Size - ini[posicion]
		posicion++
	}
	for i := 0; i < posicion; i++ {
		if disponible[i] >= c.Size {
			part.Start = ini[i]
			part.Status = getChar("1")
			mbr.Particiones[posicion] = part
			break
		}
	}
	if part.Status == getChar("0") {
		fmt.Println("Error: No hay espacio suficiente para crear la particion")
	} else {
		mbr.Particiones = ordenarParticiones(mbr.Particiones, posicion)
	}
	return mbr
}

func peorAjuste(ebr, ebraux estructuras.EBR, f *os.File, p, u *listaEBR, c Fdisk) {
	if u.size >= c.Size {
		if u.inicio == u.antes {
			f.Seek(int64(u.inicio), 0)
			ebr = leerEBR(f)
			f.Close()
			ebr.Status = getChar("1")
			ebr.Fit = getChar(c.Fit)
			ebr.Start = u.inicio
			ebr.Size = c.Size
			for i := 0; i < 16; i++ {
				ebr.Name[i] = c.Name[i]
			}
			escribirEBR(c, &ebr)
		} else {
			f.Seek(int64(u.antes), 0)
			ebr = leerEBR(f)
			f.Close()
			auxini := ebr.Next
			ebraux.Status = getChar("1")
			ebraux.Fit = getChar(c.Fit)
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
}

func ajustar(ebr, ebraux estructuras.EBR, f *os.File, p, u *listaEBR, c Fdisk) {
	var aux *listaEBR
	aux = p
	for aux != u.sig {
		if aux.size >= c.Size {
			if aux.antes == aux.inicio {
				f.Seek(int64(aux.inicio), 0)
				ebr = leerEBR(f)
				f.Close()
				ebr.Status = getChar("1")
				ebr.Fit = getChar(c.Fit)
				ebr.Start = aux.inicio
				ebr.Size = c.Size
				for i := 0; i < 16; i++ {
					ebr.Name[i] = c.Name[i]
				}
				escribirEBR(c, &ebr)
			} else {
				f.Seek(int64(aux.antes), 0)
				ebr = leerEBR(f)
				f.Close()
				auxini := ebr.Next
				ebraux.Status = getChar("1")
				ebraux.Fit = getChar(c.Fit)
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
			if parts[j].Start >= parts[j+1].Start {
				aux = parts[j]
				parts[j] = parts[j+1]
				parts[j+1] = aux
			}
		}
	}
	return parts
}

func agregarEBR(f *os.File, exte estructuras.Particion, c Fdisk, ebr estructuras.EBR) {
	var primero *listaEBR
	var ultimo *listaEBR
	var aux estructuras.EBR
	var size int
	var pos int
	if ebr.Status == getChar("0") && ebr.Next != -1 {
		f.Seek(int64(ebr.Next), 0)
		aux = leerEBR(f)
		size = aux.Start - ebr.Start
		pos = ebr.Start
		listarEBR(primero, ultimo, size, pos, pos)
		ebr = aux
	}
	for ebr.Next != -1 {
		f.Seek(int64(ebr.Next), 0)
		aux = leerEBR(f)
		pos = ebr.Size + ebr.Start
		size = aux.Start - pos
		listarEBR(primero, ultimo, size, pos, ebr.Start)
		ebr = aux
	}
	if ebr.Next == -1 {
		size = exte.Size + exte.Start - (ebr.Start + ebr.Size)
		pos = ebr.Start + ebr.Size
		listarEBR(primero, ultimo, size, pos, ebr.Start)
	}
	switch exte.Fit {
	case getChar("f"):
		ajustar(ebr, aux, f, primero, ultimo, c)
	case getChar("b"):
		ordenarEBR(primero, ultimo)
	default:
		ordenarEBR(primero, ultimo)
		peorAjuste(ebr, aux, f, primero, ultimo, c)
	}
}

func listarEBR(p, u *listaEBR, size, actual, anterior int) {
	if size > 0 {
		var nuevo *listaEBR
		nuevo.size = size
		nuevo.antes = anterior
		nuevo.inicio = actual
		if p.inicio == 0 {
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
		auxsize := 0
		auxini := 0
		auxantes := 0
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
	binary.Write(&ebrb, binary.BigEndian, &ebr)
	EscribirBytes(f, ebrb.Bytes())
	f.Close()
}

func getExtendida(mbr estructuras.MBR) estructuras.Particion {
	var part estructuras.Particion
	for i := 0; i < 4; i++ {
		if mbr.Particiones[i].Tipo == getChar("e") {
			return mbr.Particiones[i]
		}
	}
	return part
}

func comprobarParticionesVacias(mbr estructuras.MBR) bool {
	for i := 0; i < 4; i++ {
		if mbr.Particiones[i].Status != getChar("0") {
			return false
		}
	}
	return true
}

func comprobarParticionesLlenas(mbr estructuras.MBR) bool {
	for i := 0; i < 4; i++ {
		if mbr.Particiones[i].Status == getChar("1") {
			return false
		}
	}
	return true
}

func comprobarName(mbr estructuras.MBR, name string) bool {
	for i := 0; i < 4; i++ {
		n := ""
		for j := 0; j < 16; j++ {
			n += string(mbr.Particiones[i].Name[j])
		}
		if strings.EqualFold(n, name) {
			return false
		}
	}
	return true
}

func nuevaParticion(c Fdisk, mbr estructuras.MBR) estructuras.Particion {
	var part estructuras.Particion
	part.Size = c.Size
	for i := 0; i < 16; i++ {
		part.Name[i] = c.Name[i]
	}
	part.Status = getChar("1")
	part.Fit = getChar(c.Fit)
	part.Tipo = getChar(c.Tipo)
	part.Start = int(unsafe.Sizeof(mbr))
	return part
}

func nuevoEBR(c Fdisk) estructuras.EBR {
	var ebr estructuras.EBR
	ebr.Fit = byte(c.Fit[0])
	for i := 0; i < 16; i++ {
		ebr.Name[i] = byte(c.Name[i])
	}
	ebr.Next = -1
	ebr.Size = c.Size
	ebr.Status = getChar("0")
	return ebr
}

func getChar(s string) byte {
	return byte(s[0])
}

func LeerBytes(f *os.File, n int) []byte {
	bytes := make([]byte, n)
	_, err := f.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
