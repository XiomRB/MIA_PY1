package reportes

import (
	"Archivos/PY1/comandos/disco"
	"Archivos/PY1/comandos/sistema"
	"fmt"
	"strconv"
)

func RepTreeDirectory(comando Reporte) string {
	letra, indice := sistema.EncontrarMontada(comando.Id)
	if letra == -1 {
		return "Error: la particion no esta montada"
	}
	fmt.Println("Seleccione un directorio")
	nombre := ""
	var op int
	for i := 0; i < len(disco.DiscosMontados[letra].Particiones[indice].AVD); i++ {
		nombre = string(disco.DiscosMontados[letra].Particiones[indice].AVD[i].Nombre[:])
		if len(nombre) != 0 {
			fmt.Println(i+1, " ", nombre)
		}
	}
	fmt.Scanln(&op)
	if op <= len(disco.DiscosMontados[letra].Particiones[indice].AVD) {
		dot := "digraph g{\n"
		dot += graphD(&disco.DiscosMontados[letra].Particiones[indice], op-1)
		dot += "}"
		if len(dot) == 12 {
			return "La carpeta no contiene archivos"
		}
		return dot
	}
	return "Error: la opcion es incorrecta"
}

func graphD(part *disco.Montada, c int) string {
	nombre := extraerStr(part.AVD[c].Nombre)
	if part.AVD[c].IndiceDD != -1 && len(nombre) > 0 {
		dot := ""
		dot += nombre + strconv.Itoa(c) + "[shape = box];\n" + nombre + strconv.Itoa(c) + " -> detalle" + strconv.Itoa(int(part.AVD[c].IndiceDD)) + ";\n"
		dot += graphDetalle(part, int(part.AVD[c].IndiceDD), false)
		if part.AVD[c].IndiceNext != -1 {
			dot += graphD(part, int(part.AVD[c].IndiceNext))
			dot += nombre + strconv.Itoa(c) + " -> " + nombre + strconv.Itoa(int(part.AVD[c].IndiceNext)) + " [color = purple];\n"
		}
		return dot
	}
	return ""
}
