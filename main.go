package main

import (
	"Archivos/PY1/comandos"
	"fmt"
)

func main() {

	Path := "exec -path->/home/gabriela/Documentos/entrada2.mia"

	fmt.Println("Introduzca un comando:")

	comandos.Ejecutar(Path)
}
