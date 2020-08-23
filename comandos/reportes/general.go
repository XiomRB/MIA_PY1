package reportes

import (
	"fmt"
	"strings"
)

func AdministrarReportes(comando Reporte) {
	switch strings.ToLower(comando.Name) {
	case "mbr":
		fmt.Println("mbr")
	case "disk":
		fmt.Println(Reportar(comando, RepDisk(comando)))
	}
}
