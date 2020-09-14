package main

import (
	"Archivos/PY1/comandos"
	"bufio"
	"fmt"
	"os"
)

func main() {
	/*Path := "exec -path->/home/gabriela/Documentos/entrada.mia"

	comandos.Ejecutar(Path)*/
	fmt.Println("Introduzca un comando:")
	reader := bufio.NewReader(os.Stdin)
	entrada, _ := reader.ReadString('\n')
	var path string
	path = string(entrada)
	for path != "exit\n" {
		comandos.Ejecutar(path)
		fmt.Println("\nIntroduzca un comando:")
		entrada, _ := reader.ReadString('\n')
		path = string(entrada)
	}
}
