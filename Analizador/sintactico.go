package Analizador

import "fmt"

//Parser is
func Parser(cadena string) {
	tokens := Scanner(cadena)
	numpreanalisis := 0
	if len(tokens) != 0 {
		preanalisis := tokens[0]
		lc(preanalisis, tokens, numpreanalisis)
	}
}

func com(preanalisis Token, tokens []Token, numpreanalisis int) int {
	if preanalisis.tipo == 3 {
		numpreanalisis = match(3, preanalisis, numpreanalisis)
		return lp(tokens[numpreanalisis], tokens, numpreanalisis)
	} else {
		imprimirErrorSint(preanalisis, 0)
		return numpreanalisis
	}
}

func lcom(preanalisis Token, tokens []Token, numpreanalisis int) int {
	if preanalisis.tipo == 3 {
		numpreanalisis = com(preanalisis, tokens, numpreanalisis)
		return lcom(tokens[numpreanalisis], tokens, numpreanalisis)
	}
	return numpreanalisis
}

func lc(preanalisis Token, tokens []Token, numpreanalisis int) int {
	if preanalisis.tipo == 3 {
		numpreanalisis = com(preanalisis, tokens, numpreanalisis)
		return lcom(tokens[numpreanalisis], tokens, numpreanalisis)
	} else {
		imprimirErrorSint(preanalisis, 0)
		return numpreanalisis
	}
}

func param(preanalisis Token, tokens []Token, numpreanalisis int) int {
	if preanalisis.tipo == 0 {
		numpreanalisis = match(0, preanalisis, numpreanalisis)
		numpreanalisis = match(3, tokens[numpreanalisis], numpreanalisis)
		numpreanalisis = match(5, tokens[numpreanalisis], numpreanalisis)
		return p(tokens[numpreanalisis], tokens, numpreanalisis)
	} else {
		imprimirErrorSint(preanalisis, 0)
		return numpreanalisis
	}
}

func lp(preanalisis Token, tokens []Token, numpreanalisis int) int {
	if preanalisis.tipo == 0 {
		numpreanalisis = param(tokens[numpreanalisis], tokens, numpreanalisis)
		return lp(tokens[numpreanalisis], tokens, numpreanalisis)
	}
	return numpreanalisis
}

func p(preanalisis Token, tokens []Token, numpreanalisis int) int {
	switch preanalisis.tipo {
	case 1:
		return match(1, preanalisis, numpreanalisis) //numero
	case 2:
		return match(2, preanalisis, numpreanalisis) //cadena
	case 3:
		return match(3, preanalisis, numpreanalisis) //id
	case 4:
		return match(4, preanalisis, numpreanalisis) //ruta
	default:
		imprimirErrorSint(tokens[numpreanalisis], preanalisis.tipo)
		return numpreanalisis
	}
}

func match(tipo int8, preanalisis Token, numpreanalisis int) int {
	if tipo != preanalisis.tipo {
		imprimirErrorSint(preanalisis, tipo)
	} else {
		if preanalisis.tipo != 6 {
			numpreanalisis++
		}
	}
	return numpreanalisis
}

func imprimirErrorSint(t Token, tipo int8) {
	fmt.Print("Error Sintactico: Se esperaba" + obtenerTipo(tipo) + " linea ")
	fmt.Print(t.linea)
	fmt.Print(" columna ")
	fmt.Println(t.columna)
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
