package comandos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"
)

//Mkdisk info del comando mkdisk
type Mkdisk struct {
	Size int
	Path string
	Name string
	Unit string
}

//CrearDisco llamado para mkdisk
func CrearDisco(param Mkdisk) string {
	path := param.Path + param.Name
	file, err := os.Open(path)
	if err != nil {
		file, err = os.Create(path)
		if err != nil {
			if !CrearCarpeta(param.Path) {
				return "Error: No se pudo crear el directorio"
			}
			file, err = os.Create(path)
		}
		if err != nil {
			return "Error: No se pudo crear el disco"
		}
		t := time.Now()
		var disco MBR
		disco.Size = param.Size
		s := string(t.Format("Mon Jan _2 15:04:05 2006"))
		for i := 0; i < len(s); i++ {
			disco.Creacion[i] = rune(s[i])
		}
		var cero int8 = 0
		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, &cero)
		EscribirBytes(file, binario.Bytes())
		file.Seek(int64(disco.Size-1), 0)
		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, &cero)
		EscribirBytes(file, binario2.Bytes())

		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, &disco)
		EscribirBytes(file, binario3.Bytes())
		file.Close()
		return "Disco creado"
	}
	file.Close()
	return "Error: El disco ya existe"
}

func eliminarDisco(path string, l int) {
	err := os.Remove(path)
	if err != nil {
		fmt.Println("Error: El disco no existe --Linea: ", l)
	} else {
		fmt.Println("Disco borrado exitosamente")
	}
}

//CrearCarpeta funcion para crear directorios en caso no existen
func CrearCarpeta(path string) bool {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return false
	}
	return true
}

//EscribirBytes funcion para escribir en el archivo
func EscribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}
