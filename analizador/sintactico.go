package analizador

import "fmt"

var Error bool
var preanalisis Token
var numpreanalisis int
var tokens []Token

//Parser is
func Parser(cadena string) Nodo {
	tokens = Scanner(cadena)
	numpreanalisis = 0
	Error = false
	if len(tokens) != 0 {
		preanalisis = tokens[0]
		raiz := CrearNodo("RAIZ", -2, 0, 0)
		raiz.Hijos = lc()
		return raiz
	}
	return CrearNodo("", -1, 0, 0)
}

func com() Nodo {
	nuevo := CrearNodo("COMANDO", -1, 0, 0)
	if preanalisis.tipo == 3 {
		nuevo = match(3, 0) //COMANDO
		nuevo.Hijos = lp()
	} else {
		imprimirErrorSint(3)
	}
	return nuevo
}

func lcom() []Nodo {
	var lista []Nodo
	if preanalisis.tipo == 3 {
		lista = append(lista, com())
		lista1 := lcom()
		if len(lista1) != 0 {
			for i := 0; i < len(lista1); i++ {
				lista = append(lista, lista1[i])
			}
		}
	}
	return lista
}

func lc() []Nodo {
	var lista []Nodo
	if preanalisis.tipo == 3 {
		lista = append(lista, com())
		lista1 := lcom()
		if len(lista1) != 0 {
			for i := 0; i < len(lista1); i++ {
				lista = append(lista, lista1[i])
			}
		}
	} else {
		imprimirErrorSint(3)
	}
	return lista
}

func param() Nodo {
	nuevo := CrearNodo("PARAMETRO", -1, preanalisis.linea, preanalisis.columna)
	if preanalisis.tipo == 0 {
		match(0, 0)
		nuevo = match(3, 5) //5 parametro
		if preanalisis.tipo == 5 {
			match(5, 5)
			nuevo.Hijos = append(nuevo.Hijos, p())
		}
	} else {
		imprimirErrorSint(0)
	}
	return nuevo
}

func lp() []Nodo {
	var lista []Nodo
	if preanalisis.tipo == 0 {
		lista = append(lista, param())
		lista1 := lp()
		if len(lista1) != 0 {
			for i := 0; i < len(lista1); i++ {
				lista = append(lista, lista1[i])
			}
		}
	}
	return lista
}

/*func hijos()Nodo{
	nuevo := CrearNodo(preanalisis.lexema, preanalisis.tipo, preanalisis.linea, preanalisis.columna)
	if preanalisis.tipo == 5{
		match(5, 5)
		return p()
	}
	return null
}*/

func p() Nodo {
	nuevo := CrearNodo(preanalisis.lexema, preanalisis.tipo, preanalisis.linea, preanalisis.columna)
	switch preanalisis.tipo {
	case 1:
		match(1, 1) //numero
	case 2:
		match(2, 2) //cadena
	case 3:
		match(3, 3) //id
	case 4:
		match(4, 4) //ruta
	default:
		imprimirErrorSint(3)
	}
	return nuevo
}

func match(tipo, t int8) Nodo {
	nuevo := CrearNodo(preanalisis.lexema, t, preanalisis.linea, preanalisis.columna)
	if tipo != preanalisis.tipo {
		imprimirErrorSint(tipo)
	} else {
		if preanalisis.tipo != 6 {
			numpreanalisis++
			preanalisis = tokens[numpreanalisis]
		}
	}
	return nuevo
}

func imprimirErrorSint(tipo int8) {
	Error = true
	fmt.Print("Error Sintactico: Se esperaba" + obtenerTipo(tipo) + " linea ")
	fmt.Print(preanalisis.linea)
	fmt.Print(" columna ")
	fmt.Println(preanalisis.columna)
}

func obtenerTipo(t int8) string {
	switch t {
	case 0:
		return " - "
	case 1:
		return " numero "
	case 2:
		return " ruta "
	case 3:
		return " id "
	case 4:
		return " ruta "
	case 5:
		return " -> "
	default:
		return " ultimo "
	}
}
