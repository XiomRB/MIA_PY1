package disco

import (
	"Archivos/PY1/comandos"
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"unsafe"
)

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

func crearParticion(comando Fdisk) {

}

func modificarSize(comando Fdisk) {

}

func borrarPart(comando Fdisk) {

}

func escribirEnDisco() {

}

func leerDisco(path string) comandos.MBR {
	file, err := os.Open(path)
	defer file.Close()
	mbr := comandos.MBR{}
	var size int = int(unsafe.Sizeof(mbr))
	m := leerBytes(file, size)
	buffer := bytes.NewBuffer(m)
	err = binary.Read(buffer, binary.BigEndian, &mbr)
	if err != nil {
		log.Fatal(err)
	}
	return mbr
}

func LeerBytes() []byte {
	bytes := make([]byte, n)
	_, err := file.Read(bytes)
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
