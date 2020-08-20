package reportes

import (
	"log"
	"os"
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
	/*ini := "digraph G{\nrank = \"sink\"; mbr1 [ shape=plaintext label=< <table border='0' cellborder='1' cellspacing='0'> <tr> <td colspan='2'> MBR </td> </tr> <tr>  <td>Nombre</td>  <td>Valor</td> </tr>"
	mbr := disco.LeerDisco(comando.Path)*/

}

func RepDisk() {

}
