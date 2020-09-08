package reportes

import (
	"Archivos/PY1/comandos/disco"
	"Archivos/PY1/comandos/sistema"
	"Archivos/PY1/estructuras"
	"strconv"
)

func RepDirectory(comando Reporte) string {
	letra, indice := sistema.EncontrarMontada(comando.Id)
	if letra == -1 {
		return "Error, la particion no esta montada"
	}
	dot := "digraph g{\nnode [shape = box];\n"
	dot += GraphArbol(disco.DiscosMontados[letra].Particiones[indice].AVD)
	dot += "}"
	return dot
}

func GraphArbol(arbol []estructuras.AVD) string {
	dot := ""
	for i := 0; i < len(arbol); i++ {
		if arbol[i].Nombre[0] != 0 {
			dot += strconv.Itoa(i) + " [label = \"" + extraerStr(arbol[i].Nombre) + "\"]\n"
			for j := 0; j < 6; j++ {
				if arbol[i].IndicesSubs[j] != 0 {
					dot += strconv.Itoa(i) + " -> " + strconv.Itoa(int(arbol[i].IndicesSubs[j])) + "\n"
				}
			}
			if arbol[i].IndiceNext != 0 {
				dot += strconv.Itoa(i) + " -> " + strconv.Itoa(int(arbol[i].IndiceNext))
				if arbol[i].Nombre == arbol[arbol[i].IndiceNext].Nombre {
					dot += "{rank = same " + strconv.Itoa(i) + " " + strconv.Itoa(int(arbol[i].IndiceNext)) + " }\n"
				}
			}
		}
	}
	return dot
}
