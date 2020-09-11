package reportes

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

type Reporte struct {
	Name string
	Path string
	Id   string
	Ruta string
}

func AdministrarReportes(comando Reporte) {
	switch strings.ToLower(comando.Name) {
	case "mbr":
		fmt.Println(Reportar(comando, RepMBR(comando)))
	case "disk":
		fmt.Println(Reportar(comando, RepDisk(comando)))
	case "sb":
		fmt.Println(Reportar(comando, RepSB(comando)))
	case "directorio":
		fmt.Println(Reportar(comando, RepDirectory(comando)))
	case "tree_file":
		if len(comando.Ruta) == 0 {
			fmt.Println("Error: no coloco la ruta del archivo")
		} else {
			fmt.Println(Reportar(comando, RepTreeFile(comando)))
		}
	case "bm_arbdir", "bm_block", "bm_detdir", "bm_inode":
		selBitmap(comando)
	case "tree_directorio":
		fmt.Println(Reportar(comando, RepTreeDirectory(comando)))
	case "tree_complete":
		fmt.Println(Reportar(comando, repTree(comando)))
	}
}

func Reportar(comando Reporte, rep string) string {
	if len(rep) <= 35 {
		return rep
	}
	dotFile := "reporte.dot"
	pngFile := obtenerCarpeta(comando.Path)
	ioutil.WriteFile(dotFile, []byte(rep), 0755)
	png, e := exec.Command("dot", "-Tpng", dotFile).Output()
	if e != nil {
		log.Fatal(e)
	}
	ioutil.WriteFile(pngFile, png, 0755)
	return "Reporte creado"
}
