package main

import "Archivos/PY1/Analizador"

func main() {
	/*var tokens []Analizador.Token = Analizador.Scanner("mkdir -name -> arroba.dsk")
	for i := 0; i < len(tokens); i++ {
		fmt.Println(tokens[i])
	}*/
	var raiz Analizador.Nodo
	raiz = Analizador.Parser("mkdir -name -> arroba.dsk\n#comentario para probar\nmkdisk    -path ->  \"mi disco/disco 1.dsk\"\n exec")
	Analizador.ImprimirArbol(raiz)
}
