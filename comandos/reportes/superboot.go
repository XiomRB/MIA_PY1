package reportes

import (
	"Archivos/PY1/comandos/disco"
	"Archivos/PY1/comandos/sistema"
	"Archivos/PY1/estructuras"
)

func RepSB(comando Reporte) string {
	letra, indice, path := sistema.EncontrarMontada(comando.Id)
	if letra == -1 {
		return "Error, la particion no esta montada"
	}
	dot := `digraph g{
		tbl[
			shape = plaintext
			label = <
			<table color='blue' cellspacing='3'>`
	dot += GraphSB(disco.DiscosMontados[letra].Particiones[indice].Superboot)
	dot += `</table>>];}`
	return dot
}

func GraphSB(superboot estructuras.SBoot) string {
	dot := ""
	dot += "<tr><td>Nombre</td><td>" + extraerNombre(superboot.NameDisc) + "</td></tr>"
	dot += "<tr><td>Count AVD</td><td>" + Itoa(superboot.NoArbolVirtual) + "</td></tr>"
	dot += "<tr><td>Count Detalle Directorio</td><td>" + Itoa(superboot.NoDetalleDirectorio) + "</td></tr>"
	dot += "<tr><td>Count Inodos</td><td>" + Itoa(superboot.NoInodos) + "</td></tr>"
	dot += "<tr><td>Count Bloques</td><td>" + Itoa(superboot.NoBloques) + "</td></tr>"
	dot += "<tr><td>Count Free AVD</td><td>" + Itoa(superboot.NoLibreArbolVirtual) + "</td></tr>"
	dot += "<tr><td>Count Free Detalle Directorio</td><td>" + Itoa(superboot.NoLibreDetalleDirec) + "</td></tr>"
	dot += "<tr><td>Count Free Inodos</td><td>" + Itoa(superboot.NoLibreInodos) + "</td></tr>"
	dot += "<tr><td>Count Free Bloques</td><td>" + Itoa(superboot.NoLibreBloques) + "</td></tr>"
	hora := ""
	for i := 0; superboot.Creacion[i] != 0; i++ {
		hora += superboot.Creacion[i]
	}
	dot += "<tr><td>Creacion</td><td>" + hora + "</td></tr>"
	hora = ""
	for i := 0; superboot.LastMontaje[i] != 0; i++ {
		hora += superboot.LastMontaje[i]
	}
	dot += "<tr><td>Ultimo Montaje</td><td>" + hora + "</td></tr>"
	dot += "<tr><td>Inicio Bitmap AVD</td><td>" + Itoa(superboot.BitmapArbol) + "</td></tr>"
	dot += "<tr><td>Inicio AVD</td><td>" + Itoa(superboot.InicioArbol) + "</td></tr>"
	dot += "<tr><td>Inicio Bitmap Detalle Directorio</td><td>" + Itoa(superboot.BitmapDetalleDirec) + "</td></tr>"
	dot += "<tr><td>Inicio Detalle Directorio</td><td>" + Itoa(superboot.InicioDetalleDirec) + "</td></tr>"
	dot += "<tr><td>Inicio Bitmap Inodos</td><td>" + Itoa(superboot.BitmapTablaInodo) + "</td></tr>"
	dot += "<tr><td>Inicio Inodo</td><td>" + Itoa(superboot.InicioInodo) + "</td></tr>"
	dot += "<tr><td>Inicio Bitmap Bloques</td><td>" + Itoa(superboot.BitmapBloques) + "</td></tr>"
	dot += "<tr><td>Inicio Bloques</td><td>" + Itoa(superboot.InicioBloque) + "</td></tr>"
	dot += "<tr><td>Inicio Log</td><td>" + Itoa(superboot.InicioBitacora) + "</td></tr>"
	dot += "<tr><td>Size AVD</td><td>" + Itoa(superboot.SizeArbol) + "</td></tr>"
	dot += "<tr><td>Size Detalle Directorio</td><td>" + Itoa(superboot.SizeDetalleDirec) + "</td></tr>"
	dot += "<tr><td>Size Inodo</td><td>" + Itoa(superboot.SizeInodo) + "</td></tr>"
	dot += "<tr><td>Size Bloque</td><td>" + Itoa(superboot.SizeBloque) + "</td></tr>"
	dot += "<tr><td>Free Bit AVD</td><td>" + Itoa(superboot.LibreBitArbol) + "</td></tr>"
	dot += "<tr><td>Free Bit Detalle Directorio</td><td>" + Itoa(superboot.LibreBitDetalle) + "</td></tr>"
	dot += "<tr><td>Free Bit Inodos</td><td>" + Itoa(superboot.LibreBitInodo) + "</td></tr>"
	dot += "<tr><td>Free Bit Bloques</td><td>" + Itoa(superboot.LibreBitBloque) + "</td></tr>"
	dot += "<tr><td>MagicNum</td><td>" + Itoa(superboot.MagicNum) + "</td></tr>"
	return dot
}
