package main

import (
	"Archivos/PY1/Analizador"
	"fmt"
)

func main() {
	var tokens []Analizador.Token = Analizador.Scanner("der8/as/si.dsk\"/comida/t.dsk\" 888 -989 -mkdisk -> -> -988 ->9875 ->-esded")
	fmt.Println(len(tokens))
	for i := 0; i < len(tokens); i++ {
		fmt.Println(tokens[i])
	}
}
