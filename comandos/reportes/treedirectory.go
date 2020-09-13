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
		nombre = extraerStr(disco.DiscosMontados[letra].Particiones[indice].AVD[op-1].Nombre)
		if len(nombre) == 0 {
			return "Error: el directorio no existe"
		}
		if disco.DiscosMontados[letra].Particiones[indice].AVD[op-1].IndiceDD == -1 {
			return "La carpeta no contiene archivos"
		}
		dot := "digraph g{\n"
		dot += nombre + "[shape = box];\n" + nombre + " -> detalle" + strconv.Itoa(int(disco.DiscosMontados[letra].Particiones[indice].AVD[op-1].IndiceDD)) + ";\n"
		dot += graphDetalle(&disco.DiscosMontados[letra].Particiones[indice], int(disco.DiscosMontados[letra].Particiones[indice].AVD[op-1].IndiceDD), false)
		dot += "}"
		return dot
	}
	return "Error: la opcion es incorrecta"
}
