package comandos

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"os/user"
	"strings"
	"time"
)

//Mkdisk info del comando mkdisk
type Mkdisk struct {
	Size int64
	Path string
	Name string
	Unit string
}

//CrearDisco llamado para mkdisk
func CrearDisco(param Mkdisk) string {
	path := unirPath(param.Path, param.Name)
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
		file.Seek(disco.Size-1, 0)
		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, &cero)
		EscribirBytes(file, binario2.Bytes())

		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, &disco)
		EscribirBytes(file, binario3.Bytes())
		return "Disco creado"
	} else {
		file.Close()
		return "Error: El disco ya existe"
	}
}

func unirPath(p, n string) string { //-----------------------BORRAR--------------------------
	pp := ""
	if p[0] == 34 {
		r := []rune(p)
		pp = string(r[0 : len(r)-1])
		pp += n + "\""
	} else {
		pp = p + n
	}
	return pp
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

func HomePath(path string) string {
	lista := strings.Split(path, "/")
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	if lista[1] == "home" {
		path = u.HomeDir + "/"
		for i := 2; i < len(lista)-1; i++ {
			path += lista[i] + "/"
		}
	}
	return path
}

func CrearMKDISK() {

}
