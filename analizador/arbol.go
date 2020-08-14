package analizador

import "fmt"

//Nodo is
type Nodo struct {
	Tipo    int8
	Dato    string
	Linea   int16
	Columna int16
	Hijos   []Nodo
}

//CrearNodo is
func CrearNodo(v string, tipo int8, l, c int16) Nodo {
	var n Nodo
	n.Linea = l
	n.Tipo = tipo
	n.Dato = v
	n.Columna = c
	return n
}

//ImprimirArbol is
func ImprimirArbol(raiz Nodo) {
	for i := 0; i < len(raiz.Hijos); i++ {
		fmt.Println(raiz.Hijos[i])
	}
}
