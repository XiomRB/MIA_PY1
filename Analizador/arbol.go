package Analizador

import "fmt"

//Nodo is
type Nodo struct {
	tipo    int8
	dato    string
	linea   int16
	columna int16
	hijos   []Nodo
}

//CrearNodo is
func CrearNodo(v string, tipo int8, l, c int16) Nodo {
	var n Nodo
	n.linea = l
	n.tipo = tipo
	n.dato = v
	n.columna = c
	return n
}

//ImprimirArbol is
func ImprimirArbol(raiz Nodo) {
	for i := 0; i < len(raiz.hijos); i++ {
		fmt.Println(raiz.hijos[i])
	}
}
