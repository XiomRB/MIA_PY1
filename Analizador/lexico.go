package Analizador

import (
	"fmt"
	"strings"
)

//Token is
type Token struct {
	tipo    int8 // 0-  1D 2ruta 3id 4ruta 5->
	lexema  string
	linea   int16
	columna int16
}

func scanner(cadena string) []Token {
	estado := 0
	var tokens []Token
	var tok string
	var l int16 = 0 //linea
	var j int16 = 0 //columna
	for i := 0; i < len(cadena); i++ {
		switch estado {
		case 0:
			if cadena[i] == 45 { // -
				tok += string(cadena[i])
				estado = 1
			} else if cadena[i] > 47 && cadena[i] < 58 { //digito
				tok += string(cadena[i])
				estado = 2
			} else if cadena[i] == 34 { // "
				estado = 3
				tok += string(cadena[i])
			} else if (cadena[i] == 95) || (cadena[i] == 47) || (cadena[i] > 64 && cadena[i] < 91) || (cadena[i] > 96 && cadena[i] < 123) { //   letra, _ o /
				tok += string(cadena[i])
				estado = 4
			} else if cadena[i] == 10 {
				l++
				j = 0
			} else if (cadena[i] != 9) || (cadena[i] != 32) {
				fmt.Println("Error Lexico: " + string(cadena[i]) + " en la linea " + string(l) + " columna " + string(j))
			}
		case 1:
			if cadena[i] > 47 && cadena[i] < 58 {
				estado = 2
				tok += string(cadena[i])
			} else {
				tok = ""
				estado = 0
				if cadena[i] == 62 { // >
					tokens = append(tokens, crearToken("->", l, j, 5))
				} else {
					tokens = append(tokens, crearToken("-", l, j, 0))
					i--
					j--
				}
			}
		case 2:
			tok += string(cadena[i])
			if (cadena[i] == 95) || (cadena[i] == 47) || (cadena[i] > 64 && cadena[i] < 91) || (cadena[i] > 96 && cadena[i] < 123) {
				estado = 4
			} else if cadena[i] < 47 || cadena[i] > 57 { //digit{
				tokens = append(tokens, crearToken(tok, l, j, 1))
				i--
				j--
				estado = 0
				tok = ""
			}
		case 3:
			tok += string(cadena[i])
			estado = 5
			if cadena[i] == 10 || cadena[i] == 9 {
				estado = 0
				i--
				tok = ""
				fmt.Println("Error lexico: " + tok + " en la linea " + string(l) + " columna " + string(j))
			}
		case 4:
			tok += string(cadena[i])
			if cadena[i] == 46 {
				if strings.Compare("dsk", string([]byte{cadena[i+1], cadena[i+2], cadena[i+3]})) == 0 {
					tokens = append(tokens, crearToken(tok, l, j, 4))
				} else {
					tokens = append(tokens, crearToken(tok, l, j, 3))
				}
				tok = ""
				i--
				j--
				estado = 0
			} else if (cadena[i] > 0 && cadena[i] < 47) || (cadena[i] > 57 && cadena[i] < 65) || (cadena[i] > 90 && cadena[i] < 95) || (cadena[i] > 122) || cadena[i] == 96 {
				tokens = append(tokens, crearToken(tok, l, j, 3))
				tok = ""
				i--
				j--
				estado = 0
			}
		case 5:
			tok += string(cadena[i])
			if cadena[i] == 10 {
				if tok[0] != 35 {
					fmt.Println("Error Lexico: " + tok + " en la linea " + string(l) + " columna " + string(j))
				}
				tok = ""
				estado = 0
				i--
			} else if cadena[i] == 34 {
				if tok[0] == 34 {
					tokens = append(tokens, crearToken(tok, l, j, 2))
				} else {
					fmt.Println("Error Lexico: " + tok + " en la linea " + string(l) + " columna " + string(j))
				}
				tok = ""
				estado = 0
				i--
				j--
			}
		}
		j++
	}
	return tokens
}

func crearToken(lex string, lin int16, col int16, tipo int8) Token {
	var t Token
	t.lexema = lex
	t.linea = lin
	t.columna = col
	t.tipo = tipo
	return t
}
