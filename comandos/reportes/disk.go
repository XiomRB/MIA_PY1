package reportes

import (
	"unsafe"
	"Archivos/PY1/comandos/disco"
	"log"
	"os"
	"strconv"
)

type Reporte struct {
	Name string
	Path string
	Id   string
	Ruta string
}

func Reportar(string Rep) {

}

func RepMBR(comando Reporte) {
	file, err := os.Create("rep.dot")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

}

/*
func DotDisk(comando Reporte) {
	file, err := os.Create("rep.dot")
	if err != nil {
		fmt.Println("Error: No se pudo crear el reporte")
		return
	}
	dot := `digraph { tbl [	  shape=plaintext label=<  <table border='1' cellborder='0' color='blue' cellspacing='1'> <tr>`
	dot += crearGraphParts(comando)
	dot += "</tr> </table> >];}"
	var bytes []byte
	copy(bytes[:], dot)
	file.Write(bytes)
	file.Close()
}

func crearGraphParts(comando Reporte) string {
	dot := ""
	tabla := "<table border='1' cellborder='0' color='blue' cellspacing='1'> <tr> <td width = '"
	t := "' border = '0'>"
	mbr := disco.LeerDisco(comando.Path)
	if mbr.Size > 0 {
		size := mbr.Particiones[0].Start - int64(unsafe.Sizeof(mbr))
		if size > 0 {
			tam := strconv.Itoa(size / mbr.Size)
			dot += "<td width: " + tam + ">" + tabla
			dot += "<tr><td> </td></tr> <tr><td> Libre </td></tr> <tr><td> </td></tr></table></td>"
		}
		for i := 0; i < 3; i++ {
			dot += "<td>"+tabla +crearGraphP(mbr, i, -1) + "</table></td>"
			size = mbr.Particiones[i+1].Start - (mbr.Particiones[i].Start + mbr.Particiones[i].Size)
			dot += "<td>" + tabla +
		}
		size = mbr.Size - mbr.Particiones[3].Start
		dot += crearGraphP(mbr, 3, size)
		return dot
	}
	return "e"
}

func crearGraphP(mbr estructuras.MBR, indice int, size int64) string {
	p := ""
	if size > 0 {
		tam := strconv.Itoa(size / mbr.Size)
		p += "<tr><td width:'" + tam + "'> </td></tr> <tr><td> Libre </td></tr> <tr><td> </td></tr>"
	}
	tam := strconv.Itoa(mbr.Particiones[indice].Size / mbr.Size)
	p += "<tr><td width: '" + tam + "'> </td></tr> <tr><td>" + string(mbr.Particiones[indice].Name[:])
	p += ` </td></tr><tr><td> </td></tr>`
	return p
}*/

func crearGraphExt() {
	
}

func graphPart(size int64, nombre string) string {
	p := `<td width = '` + strconv.Itoa(size)
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
	ext := disco.GetExtendida(&mbr)
	if mbr.Size > 0 {
		if !disco.ComprobarParticionesVacias(&mbr) { //si hay particiones en el disco
			i := 0
			size := mbr.Particiones[0].Start - unsafe.Sizeof(mbr)
			if size > 0 {
				dot += graphPart(size*600/mbr.Size, "Libre")
			}
			for i = 0; i < 3; i++ {
				if mbr.Particiones[i+1].Status == disco.GetChar("1") {
					if mbr.Particiones[i].Tipo == disco.GetChar("e") {
						//llamar a graph extendida
					} else {
						dot += graphPart(mbr.Particiones[i].Size*600/mbr.Size, string(mbr.Particiones[i].Name[:]))
					}
					size = mbr.Particiones[i+1].Start - (mbr.Particiones[i].Start - mbr.Particiones[i].Size)
					if size > 0 {
						dot += graphPart(size*600/mbr.Size, "Libre")
					}
				} else {
					dot += graphPart(mbr.Particiones[i].Size*600/mbr.Size, string(mbr.Particiones[i].Name[:]))
					dot += graphPart((mbr.Size - mbr.Particiones[i].Start - mbr.Particiones[i].Size)*600/mbr.Size,"Libre")
					break
				}
			}

		}
	}
}

func 

type Dsk struct {
	Name string
	Size string
}
