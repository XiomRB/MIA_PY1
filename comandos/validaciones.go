package comandos

import (
	"fmt"
	"strings"
)

func ValidarPath(path string, l int) bool {
	if len(path) == 0 {
		fmt.Println("Error: El parametro path es obligatorio  --- Linea: ", l)
		return false
	}
	return true
}

func VerificarSize(size, l int) bool {
	if size > 0 {
		return true
	} else if size == -1 {
		fmt.Println("Error: El parametro size es obligatorio  --- Linea: ", l)
	}
	return false
}

func VerificarName(name string, l int) bool {
	if len(name) > 1 {
		return true
	} else if len(name) == 0 {
		fmt.Println("Error: El parametro name es obligatorio  --- Linea: ", l)
	}
	return false
}

func DarSize(s int, u string) int {
	if strings.EqualFold(u, "k") {
		return s * 1024
	} else if strings.EqualFold(string(u), "m") {
		return s * 1024 * 1024
	}
	return s
}
