package main

import "Archivos/PY1/comandos"

func main() {

	Path := "exec -path->/home/gabriela/Documentos/entrada2.mia"

	comandos.Ejecutar(Path)
	/*n := 10
	bloques := make([]estructuras.Bloque, n)

	var a [10]estructuras.Bloque
	fmt.Println(unsafe.Sizeof(bloques))
	fmt.Println(unsafe.Sizeof(a))
	fmt.Println(unsafe.Sizeof(bloques[0]))

	fmt.Println(unsafe.Sizeof(bloques[9]))

	m := disco.Montada{}
	fmt.Println(unsafe.Sizeof(m))
	fmt.Println(unsafe.Sizeof(disco.DiscosMontados))*/
}
