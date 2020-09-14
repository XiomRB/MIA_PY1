package reportes

import (
	"Archivos/PY1/comandos/disco"
	"Archivos/PY1/comandos/sistema"
	"strconv"
)

func repTree(comando Reporte) string {
	letra, indice := sistema.EncontrarMontada(comando.Id)
	if letra == -1 {
		return "Error: la particion no esta montada"
	}
	dot := "digraph g{\n"
	dot += graphDirectorio(&disco.DiscosMontados[letra].Particiones[indice])
	dot += "}"
	return dot
}

func graphDirectorio(part *disco.Montada) string {
	dot := ""
	nombreAVD := ""
	for i := 0; i < len(part.AVD); i++ {
		nombreAVD = extraerStr(part.AVD[i].Nombre)
		if len(nombreAVD) != 0 {
			dot += "carpeta" + strconv.Itoa(i) + " [shape = box label=\"" + nombreAVD + "\"];\n"
			for j := 0; j < 6; j++ {
				if part.AVD[i].IndicesSubs[j] != -1 {
					dot += "carpeta" + strconv.Itoa(i) + " -> carpeta" + strconv.Itoa(int(part.AVD[i].IndicesSubs[j])) + ";\n"
				}
			}
			if part.AVD[i].IndiceNext != -1 {
				dot += "carpeta" + strconv.Itoa(i) + " -> carpeta" + strconv.Itoa(int(part.AVD[i].IndiceNext)) + " [color=purple];\n"
			}
			if part.AVD[i].IndiceDD != -1 {
				dot += "carpeta" + strconv.Itoa(i) + " -> detalle" + strconv.Itoa(int(part.AVD[i].IndiceDD)) + ";\n"
				dot += graphDetalle(part, int(part.AVD[i].IndiceDD), true)
			}
		}
	}
	return dot
}

func graphDetalle(part *disco.Montada, det int, direc bool) string {
	dot := "detalle" + strconv.Itoa(det) + " [shape = record label = \"{Detalle " + strconv.Itoa(det)
	nombreArch := ""
	inodo := ""
	for i := 0; i < 5; i++ {
		nombreArch = extraerStr(part.DD[det].Files[i].Nombre)
		dot += " | "
		if len(nombreArch) != 0 {
			dot += "<archdet" + strconv.Itoa(det) + strconv.Itoa(i) + "> " + nombreArch
			if direc {
				if part.DD[det].Files[i].Inodo != -1 {
					inodo += "detalle" + strconv.Itoa(det) + ":archdet" + strconv.Itoa(det) + strconv.Itoa(i) + " -> inodo" + strconv.Itoa(int(part.DD[det].Files[i].Inodo)) + ";\n"
					inodo += graphInodo(part, int(part.DD[det].Files[i].Inodo))
				}
			}
		}

	}
	if part.DD[det].Next != -1 {
		dot += " | <indirect>Detalle" + strconv.Itoa(int(part.DD[det].Next)) + "}\"];\n"
		dot += "detalle" + strconv.Itoa(det) + ":indirect -> " + "detalle" + strconv.Itoa(int(part.DD[det].Next)) + " [color=green];\n"
		dot += inodo
		dot += graphDetalle(part, int(part.DD[det].Next), direc)
	} else {
		dot += "}\"];\n"
		dot += inodo
	}
	return dot
}

func graphInodo(part *disco.Montada, ind int) string {
	dot := "inodo" + strconv.Itoa(ind) + " [shape = record label = \"{<ind" + strconv.Itoa(ind) + ">Inodo " + strconv.Itoa(ind)
	contenido := ""
	for i := 0; i < 4; i++ {
		dot += " | "
		if part.Inodos[ind].Bloques[i] != -1 {
			dot += " <inbb" + strconv.Itoa(int(part.Inodos[ind].Bloques[i])) + ">Bloque " + strconv.Itoa(int(part.Inodos[ind].Bloques[i]))
			contenido += "bb" + strconv.Itoa(int(part.Inodos[ind].Bloques[i])) + " [ shape = box label = \"" + retornarBloque(part.BB[part.Inodos[ind].Bloques[i]].Text) + "\"];\n"
			contenido += "inodo" + strconv.Itoa(ind) + ":inbb" + strconv.Itoa(int(part.Inodos[ind].Bloques[i])) + " -> " + "bb" + strconv.Itoa(int(part.Inodos[ind].Bloques[i])) + ";\n"
		}
	}
	if part.Inodos[ind].Indirecto != -1 {
		dot += " | <indirecto" + strconv.Itoa(int(part.Inodos[ind].Indirecto)) + ">Indirecto " + strconv.Itoa(int(part.Inodos[ind].Indirecto)) + "}\"];\n"
		contenido += "inodo" + strconv.Itoa(ind) + ":indirecto" + strconv.Itoa(int(part.Inodos[ind].Indirecto)) + " -> inodo" + strconv.Itoa(int(part.Inodos[ind].Indirecto)) + ":ind" + strconv.Itoa(int(part.Inodos[ind].Indirecto)) + " [color=blue];\n"
		dot += contenido
		dot += graphInodo(part, int(part.Inodos[ind].Indirecto))
	} else {
		dot += "}\"];\n"
		dot += contenido
	}
	return dot
}
