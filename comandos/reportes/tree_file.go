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
		return "Error: no existe el archivo"
	}
	dot := "digraph g{\n"
	dot += GraphFile(comando.Ruta, inodos, dd, file, &disco.DiscosMontados[letra].Particiones[indice])
	dot += "}"
	return dot
}

func GraphFile(path string, inodos []int, dd, file int, part *disco.Montada) string {
	dot, nodo, arch := GraphCarpeta(path)
	dot += "detalle" + strconv.Itoa(dd) + " [shape = record label= \"{Detalle Directorio|"
	if len(inodos) == 0 {
		dot += arch + "}\"];\n"
		return dot
	}
	dot += "<inodo" + strconv.Itoa(inodos[0]) + ">" + arch + "}\"];\n"
	dot += nodo + " -> " + "detalle" + strconv.Itoa(dd) + ";\n"
	dot += "detalle" + strconv.Itoa(dd) + " -> inodo" + strconv.Itoa(inodos[0]) + ":ind" + strconv.Itoa(inodos[0]) + ";\n"
	nombre := ""

	for i := 0; i < len(inodos); i++ {
		dot += "inodo" + strconv.Itoa(inodos[i]) + " [shape = record label = \"{<ind" + strconv.Itoa(inodos[i]) + ">Inodo " + strconv.Itoa(inodos[i])
		for j := 0; j < 4; j++ {
			if part.Inodos[inodos[i]].Bloques[j] != -1 {
				dot += " | <inbb" + strconv.Itoa(int(part.Inodos[inodos[i]].Bloques[j])) + ">Bloque " + strconv.Itoa(int(part.Inodos[inodos[i]].Bloques[j]))
				nombre += "bb" + strconv.Itoa(int(part.Inodos[inodos[i]].Bloques[j])) + " [ shape = box label = \"" + retornarBloque(part.BB[part.Inodos[inodos[i]].Bloques[j]].Text) + "\"];\n"
				nombre += "inodo" + strconv.Itoa(inodos[i]) + ":inbb" + strconv.Itoa(int(part.Inodos[inodos[i]].Bloques[j])) + " -> " + "bb" + strconv.Itoa(int(part.Inodos[inodos[i]].Bloques[j])) + ";\n"
			}
		}
		if i < len(inodos)-1 {
			dot += " | <indirecto" + strconv.Itoa(inodos[i+1]) + ">Indirecto " + strconv.Itoa(inodos[i+1])
			nombre += "inodo" + strconv.Itoa(inodos[i]) + ":indirecto" + strconv.Itoa(inodos[i+1]) + " -> inodo" + strconv.Itoa(inodos[i+1]) + ":ind" + strconv.Itoa(inodos[i+1]) + ";\n"
		}
		dot += "}\"];\n"
	}
	dot += nombre
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
	return dot, "carpeta" + strconv.Itoa(len(ruta)-2), ruta[len(ruta)-1]
}

func retornarBloque(b [25]byte) string {
	bb := ""
	for i := 0; i < 25; i++ {
		if b[i] != 0 {
			bb += string(b[i])
		}
	}
	return bb
}
