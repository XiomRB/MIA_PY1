package main

import (
	"Archivos/PY1/comandos"
	"fmt"
)

func main() {
	//Path := "exec -path->/home/gabriela/Documentos/entrada2.mia"

	fmt.Println("Introduzca un comando:")
	var cmd string
	fmt.Scanln(&cmd)
	for len(cmd) > 4 {
		comandos.Ejecutar(cmd)
		fmt.Println("\nIntroduzca un comando:")
		fmt.Scanln(&cmd)
	}
}
