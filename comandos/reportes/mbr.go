package reportes

import (
	"Archivos/PY1/comandos/disco"
	"Archivos/PY1/estructuras"
	"strconv"
)

func RepMBR(comando Reporte) string {
	letra := comando.Id[2]
	indice := comando.Id[3]
	dot := `digraph g{
			tbl[
				shape = plaintext
				label = <
				<table color='blue' cellspacing='3'>`
	if verificarPM(letra, indice) {
		dot += graphMBR(letra) + "}"
		return dot
	}
	return "Error, la particion no esta montada"
}

func graphMBR(letra byte) string {
	dot := ""
	dot += `<tr><td colspan = '2'>MBR</td></tr>`
	mbr := disco.LeerDisco(disco.DiscosMontados[disco.EncontrarLetra(letra)].Path)
	dot += "<tr><td>Indice</td><td>" + intString(mbr.Indice) + "</td></tr>"
	n := ""
	for i := 0; mbr.Creacion[i] != 0; i++ {
		n += string(mbr.Creacion[i])
	}
	dot += "<tr><td>Time</td><td>" + n + "</td></tr>"
	dot += "<tr><td>Size</td><td>" + intString(mbr.Size) + "</td></tr>"
	for i := 0; i < 4; i++ {
		if mbr.Particiones[i].Name[0] != 0 {
			dot += "<tr><td>Particion</td><td>" + extraerNombre(mbr.Particiones[i].Name) + "</td></tr>"
		}
	}
	dot += `</table>>];`
	dot += graphParticion(disco.DiscosMontados[disco.EncontrarLetra(letra)].Path, mbr.Particiones)
	return dot
}

func graphParticion(path string, parts [4]estructuras.Particion) string {
	dot := ""
	for i := 0; i < 4; i++ {
		if parts[i].Status == disco.GetChar("1") {
			dot += "tbl" + intString(int64(i)) + `[
					shape = plaintext
					label = <<table  color='blue' cellspacing='3'>`
			dot += "<tr><td colspan = '2'> Particion " + extraerNombre(parts[i].Name) + "</td></tr>"
			dot += "<tr><td>Status</td><td>" + string(parts[i].Status) + "</td></tr>"
			dot += "<tr><td>Tipo</td><td>" + string(parts[i].Tipo) + "</td></tr>"
			dot += "<tr><td>Fit</td><td>" + string(parts[i].Fit) + "</td></tr>"
			dot += "<tr><td>Start</td><td>" + intString(parts[i].Start) + "</td></tr>"
			dot += "<tr><td>Size</td><td>" + intString(parts[i].Size) + "</td></tr>"
			dot += `</table>>];`
			if parts[i].Tipo == disco.GetChar("e") {
				dot += graphEBR(path, parts[i])
			}
		}
	}
	return dot
}

func graphEBR(path string, ext estructuras.Particion) string {
	dot := ""
	ebr := disco.LeerEBR(path, ext.Start)
	for i := 0; ebr.Next != -1; i++ {
		dot += "tble" + intString(int64(i)) + `[
				shape = plaintext
				label = <<table color='blue' cellspacing='3'>`
		dot += "<tr><td colspan = '2'> EBR " + extraerNombre(ebr.Name) + "</td></tr>"
		dot += "<tr><td>Status</td><td>" + string(ebr.Status) + "</td></tr>"
		dot += "<tr><td>Fit</td><td>" + string(ebr.Fit) + "</td></tr>"
		dot += "<tr><td>Start</td><td>" + intString(ebr.Start) + "</td></tr>"
		dot += "<tr><td>Size</td><td>" + intString(ebr.Size) + "</td></tr>"
		dot += "<tr><td>Next</td><td>" + intString(ebr.Next) + "</td></tr>"
		dot += `</table>>];`
		ebr = disco.LeerEBR(path, ebr.Next)
	}
	dot += `tblef[
			shape = plaintext
			label = <<table color='blue' cellspacing='3'>`
	dot += "<tr><td colspan = '2'> EBR " + extraerNombre(ebr.Name) + "</td></tr>"
	dot += "<tr><td>Status</td><td>" + string(ebr.Status) + "</td></tr>"
	dot += "<tr><td>Fit</td><td>" + string(ebr.Fit) + "</td></tr>"
	dot += "<tr><td>Start</td><td>" + intString(ebr.Start) + "</td></tr>"
	dot += "<tr><td>Size</td><td>" + intString(ebr.Size) + "</td></tr>"
	dot += "<tr><td>Next</td><td>" + intString(ebr.Next) + "</td></tr>"
	dot += `</table>>];`
	return dot
}

func extraerNombre(nombre [16]byte) string {
	n := ""
	for i := 0; i < 16; i++ {
		if nombre[i] != 0 {
			n += string(nombre[i])
		}
	}
	return n
}

func extraerStr(nombre [20]byte) string {
	n := ""
	for i := 0; i < 20; i++ {
		if nombre[i] != 0 {
			n += string(nombre[i])
		}
	}
	return n
}

func intString(i int64) string {
	num := strconv.Itoa(int(i))
	return num
}
