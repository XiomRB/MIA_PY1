package main

import (
	"fmt"
	"strings"
)

func main() {
	/*
		var raiz analizador.Nodo
		raiz = analizador.Parser("mkdir -name -> arroba.dsk\n#comentario para probar\nmkdisk /*\n   -path ->  \"mi disco/disco 1.dsk\"\n exec")
		analizador.ImprimirArbol(raiz)*/
	//var disco comandos.MBR
	/*var mkd comandos.Mk
	mkd.Path = "\"/home/gabriela/Documentos/gabii/g 1/\""
	mkd.Name = "disco.dsk"
	mkd.Size = 1000
	fmt.Println(comandos.CrearDisco(mkd))

	//mt.Println(len(strings.Split(mkd.Path, "/")))*/
	fmt.Println(strings.EqualFold("HoLA", "hOla"))
}
