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
				full, indice := comprobarParticionesLlenas(mbr)
				if !comprobarName(mbr, comando.Name) {
					if comprobarParticionesVacias(mbr) {
						if strings.EqualFold(comando.Tipo, "p") { //primaria
							mbr.Particiones[0] = nuevaParticion(comando, mbr)
						} else if strings.EqualFold(comando.Tipo, "e") { //extendida
							mbr.Particiones[0] = nuevaParticion(comando, mbr)
							ebr := nuevoEBR(comando)
							ebr.Start = mbr.Particiones[0].Start
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
						} else { //logica
							fmt.Println("Error: No existe particion extendida donde pueda agregar una logica")
							return
						}
					} else if full { // 4 particiones llenas
						if strings.EqualFold(comando.Tipo, "l") {
							extendida := getExtendida(mbr)
							if extendida.Size > 0 {
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
									agregarEBR(extendida, comando, ebr)
								}
							} else {
								fmt.Println("Error: no se puede crear una particion logica, no existe una extendida")
							}
						} else {
							fmt.Println("Error: No se puede crear la particion, ha llegado al numero maximo de particiones")
						}
					} else {

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

func LeerBytes(f *os.File, n int) []byte {
	bytes := make([]byte, n)
	_, err := f.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func primerAjuste() {

}

func peorAjuste() {

}

func mejorAjuste() {

}

func comprobarName(mbr estructuras.MBR, name string) bool {
	for i := 0; i < 4; i++ {
		n := ""
		for j := 0; j < 16; j++ {
			n += string(mbr.Particiones[i].Name[j])
		}
		if strings.EqualFold(n, name) {
			return true
		}
	}
	return false
}

func comprobarParticionesVacias(mbr estructuras.MBR) bool {
	for i := 0; i < 4; i++ {
		if mbr.Particiones[i].Size != 0 {
			return false
		}
	}
	return true
}

func comprobarParticionesLlenas(mbr estructuras.MBR) (bool, int) {
	for i := 0; i < 4; i++ {
		if mbr.Particiones[i].Size == 0 {
			return false, i
		}
	}
	return true, -1
}

func nuevaParticion(c Fdisk, mbr estructuras.MBR) estructuras.Particion {
	var part estructuras.Particion
	part.Size = c.Size
	for i := 0; i < 16; i++ {
		part.Name[i] = byte(c.Name[i])
	}
	st := "1"
	part.Status = byte(st[0])
	part.Fit = byte(strings.ToLower(c.Fit)[0])
	part.Tipo = byte(strings.ToLower(c.Tipo)[0])
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
	st := "0"
	ebr.Status = byte(st[0])
	return ebr
}

func agregarEBR(exte estructuras.Particion, c Fdisk, ebr estructuras.EBR) {

}

func getExtendida(mbr estructuras.MBR) estructuras.Particion {
	t := "e"
	var part estructuras.Particion
	for i := 0; i < 4; i++ {
		if mbr.Particiones[i].Tipo == byte(t[0]) {
			return mbr.Particiones[i]
		}
	}
	return part
}

func getChar(s string) byte {
	return byte(s[0])
}
