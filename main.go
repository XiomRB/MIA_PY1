package main

import "Archivos/PY1/comandos"

func main() {
	/*
		var raiz analizador.Nodo
		raiz = analizador.Parser("mkdir -name -> arroba.dsk\n#comentario para probar\nmkdisk /*\n   -path ->  \"mi disco/disco 1.dsk\"\n exec")
		analizador.ImprimirArbol(raiz)*/
	//var disco comandos.MBR
	Path := "exec -path->/home/gabriela/Documentos/entrada.mia"

	comandos.Ejecutar(Path)

}
