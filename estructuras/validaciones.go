package estructuras

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

func VerificarSize(size int64, l int) bool {
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

func DarSize(s int64, u string) int64 {
	if strings.EqualFold(u, "k") {
		return int64(s * 1024)
	} else if strings.EqualFold(string(u), "m") {
		return int64(s * 1024 * 1024)
	}
	return int64(s)
}
