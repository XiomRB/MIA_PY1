package main

import (
	"fmt"
)

func main() {
	saludo := `	 
	`
	fmt.Println(saludo[0])
	fmt.Println(saludo[1])
	fmt.Println(saludo[2])
	fin := ""
	if saludo[0] == 76 {
		fin = fin + string(saludo[0])
	} else {
		fin = "no coincidio"
	}

	fmt.Println(fin)
}
