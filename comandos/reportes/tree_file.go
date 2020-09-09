package reportes

import (
	"Archivos/PY1/comandos/disco"
	"Archivos/PY1/comandos/sistema"
	"strconv"
	"strings"
)

func RepTreeFile(comando Reporte) string {
	letra, indice := sistema.EncontrarMontada(comando.Id)
	if letra == -1 {
		return "Error: la particion no esta montada"
	}
	carpeta, dd, file, inodos := sistema.BuscarArchivo(&disco.DiscosMontados[letra].Particiones[indice], comando.Ruta)
	if carpeta == -1 {
		return "Error: no existe la ruta"
	}
	dot := "digraph g{\n"
	dot += GraphFile(comando.Ruta, inodos, dd, file, &disco.DiscosMontados[letra].Particiones[indice])
	dot += "}"
	return dot
}

func GraphFile(path string, inodos []int, dd, file int, part *disco.Montada) string {
	dot, nodo, arch := GraphCarpeta(path)
	dot += "detalle" + strconv.Itoa(dd) + " [shape = textplain label <<table color='blue' cellspacing='0'><tr><td colspan ='2'>Detalle Directorio</td></tr><tr><td>" + arch + "</td><td>Inodo "
	if len(inodos) == 0 {
		dot += "-1</td></tr></table>>];\n"
		return dot
	}
	dot += strconv.Itoa(inodos[0]) + "</td></tr></table>>];\n"
	dot += nodo + " -> inodo" + strconv.Itoa(inodos[0]) + ";\n"
	subinodo := ""
	nombre := ""
	for i := 0; i < len(inodos); i++ {
		nombre = "inodo" + strconv.Itoa(inodos[i])
		subinodo = "subgraph " + nombre + "{\nnode [shape = box];\n"
		for j := 0; j < 4; j++ {
			if part.Inodos[i].Bloques[j] != -1 {
				subinodo += nombre + "bb" + strconv.Itoa(j) + " [label = \"Bloque " + strconv.Itoa(int(part.Inodos[i].Bloques[j])) + "\"];\n"
				dot += "bb" + strconv.Itoa(int(part.Inodos[i].Bloques[j])) + " [label = \"" + retornarBloque(part.BB[j].Text) + "\"];\n"
				dot += nombre + "bb" + strconv.Itoa(j) + " -> " + "bb" + strconv.Itoa(int(part.Inodos[i].Bloques[j])) + ";\n"
			}
		}
		if i < len(inodos)-1 {
			subinodo += nombre + "indirecto" + " [label=\"Inodo " + strconv.Itoa(inodos[i+1]) + "\"];\n"
			dot += nombre + "indirecto -> inodo" + strconv.Itoa(inodos[i+1]) + ";\n"
		}
		subinodo += "}\n"
		dot += subinodo
	}
	return dot
}

func GraphCarpeta(path string) (string, string, string) { //retorna el dot, y el nodo de la ultima carpeta y el nombre del archivo
	path = sistema.ElimComillas(path)
	ruta := strings.Split(path, "/")
	ruta[0] = "/"
	dot := ""
	for i := 0; i < len(ruta)-1; i++ {
		dot += "carpeta" + strconv.Itoa(i) + " [shape = box label = \"" + ruta[i] + "\"];\n"
		if i < len(ruta)-2 {
			dot += "carpeta" + strconv.Itoa(i) + " -> carpeta" + strconv.Itoa(i+1) + ";\n"
		}
	}
	return dot, "carpeta" + strconv.Itoa(len(ruta)-1), ruta[len(ruta)-1]
}

func retornarBloque(b [25]byte) string {
	bb := ""
	for i := 0; b[i] != 0; i++ {
		bb += string(b[i])
	}
	return bb
}
