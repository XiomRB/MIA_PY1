package reportes

import (
	"Archivos/PY1/comandos/disco"
	"Archivos/PY1/estructuras"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

func graphParticiones(path string) string {
	dot := ""
	mbr := disco.LeerDisco(path)
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
						dot += crearGraphExt(mbr.Particiones[i], path, mbr.Size)
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

func RepDisk(comando Reporte) string {
	letra := comando.Id[2]
	indice := comando.Id[3]
	dot := `digraph g{
			tbl[
				shape = plaintext
				label = <
				<table border='1' cellborder='0' color='blue' cellspacing='1'>
        			<tr>`
	if verificarPM(letra, indice) {
		dot += graphParticiones(disco.DiscosMontados[disco.EncontrarLetra(letra)].Path)
		dot += "</tr></table>>];}"
		return dot
	}
	return "Error, la particion no esta montada"
}

func Reportar(comando Reporte, rep string) string {
	if len(rep) == 35 {
		return rep
	}
	err := ioutil.WriteFile("reporte.dot", []byte(rep), 0755)
	if err != nil {
		return "Error: No se pudo crear el reporte"
	}
	/*p := obtenerCarpeta(comando.Path)
	path, _ := exec.LookPath("dot")
	cmd, _ := exec.Command(path, "-Tpng", "reporte.dot").Output()
	mode := int(0777)
	ioutil.WriteFile(p, cmd, os.FileMode(mode))*/
	return "Reporte creado"
}

func verificarPM(letra, indice byte) bool {
	num, err := strconv.Atoi(string(indice))
	if err != nil {
		log.Fatal(err)
	}
	l := disco.EncontrarLetra(letra)
	if len(disco.DiscosMontados) > l {
		if disco.DiscosMontados[l].Estado {
			if len(disco.DiscosMontados[l].Particiones) >= num {
				return disco.DiscosMontados[l].Particiones[num-1].Estado
			}
		}
	}
	return false
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
		} else {
			path = u.HomeDir + "/"
			for i := 3; i < len(lista)-1; i++ {
				path += lista[i] + "/"
			}
		}
	} else {
		path = "/"
		for i := 0; i < len(lista)-1; i++ {
			path += lista[i] + "/"
		}
	}
	fmt.Println(path)
	os.MkdirAll(path, 0755)
	return ruta
}
