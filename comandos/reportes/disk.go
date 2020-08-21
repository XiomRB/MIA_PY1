package reportes

import (
	"Archivos/PY1/comandos/disco"
	"Archivos/PY1/estructuras"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"unsafe"
)

type Reporte struct {
	Name string
	Path string
	Id   string
	Ruta string
}

func RepMBR(comando Reporte) {
	file, err := os.Create("rep.dot")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

}

func crearGraphExt(ext estructuras.Particion, path string, size int64) string {
	p := `<td width = '` + strconv.Itoa(int(ext.Size*600/size))
	p += `' border = '0'>
			<table border='1' cellborder='0' color='blue' cellspacing='1'>
		 		<tr><td>` + string(ext.Name[:]) + `</td></tr>
				 <tr>`

	ebr := disco.LeerEBR(path, ext.Start)
	libre := ebr.Start - ext.Start
	if libre > 0 {
		p += graphPart(libre*600/size, "Libre")
	}
	aux := ebr
	for ebr.Next != -1 {
		aux = ebr
		ebr = disco.LeerEBR(path, aux.Next)
		libre := ebr.Start - aux.Start - aux.Size
		if libre > 0 {
			p += graphPart(libre*600/size, "Libre")
		}
		graphPart(ebr.Size*600/size, string(ebr.Name[:]))
	}
	libre = ext.Size - ebr.Size - ebr.Start
	if libre > 0 {
		p += graphPart(libre*600/size, "Libre")
	}
	p += `</tr>
			</table>
		</td>`
	return p
}

func graphPart(size int64, nombre string) string {
	p := `<td width = '` + strconv.Itoa(int(size))
	p += `' border = '0'>
			<table border='1' cellborder='0' color='blue' cellspacing='1'>
				<tr><td>  </td></tr>
		 		<tr><td>` + nombre + `</td></tr>
		 		<tr><td>  </td></tr>
			</table>
		</td>`
	return p
}

func graphParticiones(comando Reporte) string {
	dot := ""
	mbr := disco.LeerDisco(comando.Path)
	if mbr.Size > 0 {
		if !disco.ComprobarParticionesVacias(&mbr) { //si hay particiones en el disco
			i := 0
			size := mbr.Particiones[0].Start - int64(unsafe.Sizeof(mbr))
			if size > 0 {
				dot += graphPart(size*600/mbr.Size, "Libre")
			}
			for i = 0; i < 3; i++ {
				if mbr.Particiones[i+1].Status == disco.GetChar("1") {
					if mbr.Particiones[i].Tipo == disco.GetChar("e") {
						dot += crearGraphExt(mbr.Particiones[i], comando.Path, mbr.Size)
					} else {
						dot += graphPart(mbr.Particiones[i].Size*600/mbr.Size, string(mbr.Particiones[i].Name[:]))
					}
					size = mbr.Particiones[i+1].Start - (mbr.Particiones[i].Start - mbr.Particiones[i].Size)
					if size > 0 {
						dot += graphPart(size*600/mbr.Size, "Libre")
					}
				} else {
					dot += graphPart(mbr.Particiones[i].Size*600/mbr.Size, string(mbr.Particiones[i].Name[:]))
					dot += graphPart((mbr.Size-mbr.Particiones[i].Start-mbr.Particiones[i].Size)*600/mbr.Size, "Libre")
					break
				}
			}
			if i == 3 {
				dot += graphPart(mbr.Particiones[i].Size*600/mbr.Size, string(mbr.Particiones[i].Name[:]))
				size = mbr.Size - mbr.Particiones[3].Size - mbr.Particiones[3].Start
				if size > 0 {
					dot += graphPart(size*600/mbr.Size, "Libre")
				}
			}
		}
	}
	return dot
}

func RepDisk(comando Reporte) {

}

func Reportar(comando Reporte, rep string) string {
	file, err := os.Create("reporte.dot")
	if err != nil {
		return "Error: No se pudo crear el reporte"
	}
	var dot []byte
	copy(dot[:], rep)
	file.Write(dot)
	file.Close()
	p := obtenerCarpeta(comando.Path)
	path, _ := exec.LookPath("dot")
	cmd, _ := exec.Command(path, "-Tpng", "reporte.dot").Output()
	mode := int(0777)
	ioutil.WriteFile(p, cmd, os.FileMode(mode))
	return "Reporte creado"
}

func verificarPM() {

}

func quitarEspacios(p string) string {
	path := ""
	if p[0] == 34 {
		for i := 1; i < len(p)-1; i++ {
			if p[i] == 32 {
				path += "_"
			} else {
				path += string(p[i])
			}
		}
		return path
	}
	return p
}

func obtenerCarpeta(p string) string {
	path := quitarEspacios(p)
	ruta := path
	lista := strings.Split(path, "/")
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	if lista[1] == "home" {
		if lista[2] != u.Username {
			path = u.HomeDir + "/"
			for i := 2; i < len(lista)-1; i++ {
				path += lista[i] + "/"
			}
		}
	}
	os.MkdirAll(path, 0755)
	return ruta
}
