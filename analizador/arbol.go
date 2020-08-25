package analizador

import (
	"fmt"
	"log"
	"os/user"
	"strconv"
	"strings"
)

//Nodo is
type Nodo struct {
	Tipo    int8
	Dato    string
	Linea   int
	Columna int
	Hijos   []Nodo
}

//CrearNodo is
func CrearNodo(v string, tipo int8, l, c int) Nodo {
	var n Nodo
	n.Linea = l
	n.Tipo = tipo
	n.Dato = v
	n.Columna = c
	return n
}

//Imprimir comando
func Imprimir(raiz Nodo) {
	fmt.Print(raiz.Dato)
	for i := 0; i < len(raiz.Hijos); i++ {
		fmt.Print(raiz.Hijos[i].Dato, "->", raiz.Hijos[i].Hijos[0].Dato)
	}
	fmt.Println("")
}

func HomePath(p Nodo) string {
	path := nuevoPath(p)
	lista := strings.Split(path, "/")
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	if lista[1] == "home" {
		if lista[2] != u.Username {
			path = u.HomeDir + "/"
			for i := 2; i < len(lista)-1; i++ {
				path += lista[i] + "/"
			}
			path += lista[len(lista)-1]
		}
	}
	return path
}

func nuevoPath(p Nodo) string {
	path := ""
	if p.Dato[0] == 34 {
		for i := 1; i < len(p.Dato)-1; i++ {
			if p.Dato[i] == 32 {
				path += "_"
			} else {
				path += string(p.Dato[i])
			}
		}
		return path
	}
	return p.Dato
}

func ValidarSize(f Nodo) int64 {
	s, err := strconv.Atoi(f.Dato)
	if err == nil {
		if s > 0 {
			return int64(s)
		}
		fmt.Println("Error: El parametro size debe ser mayor a 0  --- Linea: ", f.Linea, " Col: ", f.Columna)
		return 0
	}
	fmt.Println("Error: El parametro size solo reconoce numeros  --- Linea: ", f.Linea, " Col: ", f.Columna)
	return 0
}

func ValidarFit(f Nodo) string {
	if strings.EqualFold(f.Dato, "FF") || strings.EqualFold(f.Dato, "WF") || strings.EqualFold(f.Dato, "BF") {
		return string(f.Dato[0])
	}
	fmt.Println("Error: El parametro fit no reconoce ese valor --- Linea: ", f.Linea, " Col: ", f.Columna)
	return ""
}

func ValidarUnidad(part bool, u Nodo) string {
	if strings.EqualFold(u.Dato, "k") || strings.EqualFold(u.Dato, "m") {
		return string(u.Dato[0])
	}
	if strings.EqualFold(u.Dato, "b") && part {
		return string(u.Dato[0])
	}
	fmt.Println("Error: El parametro unit no reconoce ese valor --- Linea: ", u.Linea, " Col: ", u.Columna)
	return ""
}
