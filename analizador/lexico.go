package analizador

import (
	"fmt"
	"strings"
)

//Token is
type Token struct {
	tipo    int8 // 0-  1D 2ruta 3id 4ruta 5->   6 EOF
	lexema  string
	linea   int
	columna int
}

//Scanner is
func Scanner(cadena string) []Token {
	estado := 0
	var tokens []Token
	var tok string
	cadena += "?"
	l := 0 //linea
	j := 0 //columna
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
			} else if (cadena[i] == 95) || (cadena[i] > 64 && cadena[i] < 91) || (cadena[i] > 96 && cadena[i] < 123) { //   letra, _ o /
				tok += string(cadena[i])
				estado = 4
			} else if cadena[i] == 47 { // /
				estado = 6
				tok += string(cadena[i])
			} else if cadena[i] == 92 { // \
				estado = 7
			} else if cadena[i] == 35 { //#
				estado = 8
			} else if cadena[i] == 10 {
				l++
				j = 0
			} else if cadena[i] != 9 {
				if cadena[i] != 32 {
					if (i == len(cadena)-1) && (cadena[i] == 63) {
						tokens = append(tokens, crearToken(string(cadena[i]), l, j, 6)) //EOF
					} else {
						imprimirError(string(cadena[i]), l, j)
					}
				}
			}
		case 1:
			if cadena[i] > 47 && cadena[i] < 58 { //digit
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
			if (cadena[i] == 95) || (cadena[i] == 47) || (cadena[i] > 64 && cadena[i] < 91) || (cadena[i] > 96 && cadena[i] < 123) {
				tok += string(cadena[i])
				estado = 4
			} else if cadena[i] > 47 && cadena[i] < 58 { //digit{
				tok += string(cadena[i])
			} else {
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
				imprimirError(tok, l, j)
				estado = 0
				i--
				tok = ""
			}
		case 4:
			if cadena[i] == 46 {
				if strings.Compare("dsk", string([]byte{cadena[i+1], cadena[i+2], cadena[i+3]})) == 0 {
					tokens = append(tokens, crearToken(tok+".dsk", l, j, 4))
					i += 4
					j += 4
				} else if strings.Compare("mia", string([]byte{cadena[i+1], cadena[i+2], cadena[i+3]})) == 0 {
					tokens = append(tokens, crearToken(tok+".mia", l, j, 4))
					i += 4
					j += 4
				} else if strings.Compare("png", string([]byte{cadena[i+1], cadena[i+2], cadena[i+3]})) == 0 || strings.Compare("jpg", string([]byte{cadena[i+1], cadena[i+2], cadena[i+3]})) == 0 || strings.Compare("pdf", string([]byte{cadena[i+1], cadena[i+2], cadena[i+3]})) == 0 {
					tokens = append(tokens, crearToken(tok+".png", l, j, 4))
					i += 4
					j += 4
				} else if strings.Compare("txt", string([]byte{cadena[i+1], cadena[i+2], cadena[i+3]})) == 0 {
					tokens = append(tokens, crearToken(tok+".txt", l, j, 4))
					i += 4
					j += 4
				} else {
					tokens = append(tokens, crearToken(tok, l, j, 3))
				}
				tok = ""
				i--
				j--
				estado = 0
			} else if (cadena[i] == 95) || (cadena[i] > 46 && cadena[i] < 58) || (cadena[i] > 64 && cadena[i] < 91) || (cadena[i] > 96 && cadena[i] < 123) { //   letra, _ o /
				tok += string(cadena[i])
			} else {
				tokens = append(tokens, crearToken(tok, l, j, 3))
				tok = ""
				i--
				j--
				estado = 0
			}
		case 5:
			if cadena[i] == 10 {
				if tok[0] != 35 {
					imprimirError(tok, l, j)
				}
				tok = ""
				estado = 0
				i--
			} else if cadena[i] == 34 {
				if tok[0] == 34 {
					tokens = append(tokens, crearToken(tok+"\"", l, j, 2))
				} else {
					imprimirError(tok, l, j)
				}
				tok = ""
				estado = 0
			} else {
				tok += string(cadena[i])
			}
		case 6:
			if (cadena[i] == 95) || (cadena[i] > 64 && cadena[i] < 91) || (cadena[i] > 96 && cadena[i] < 123) || (cadena[i] > 47 && cadena[i] < 58) {
				tok += string(cadena[i])
				estado = 4
			} else {
				imprimirError(tok, l, j)
			}
		case 7:
			if cadena[i] == 42 { // *
				estado = 0
				tok = ""
			} else {
				imprimirError("\\", l, j)
			}
		case 8:
			if cadena[i] == 10 {
				estado = 0
			}
		}
		j++
	}
	return tokens
}

func crearToken(lex string, lin int, col int, tipo int8) Token {
	var t Token
	t.lexema = lex
	t.linea = lin
	t.columna = col
	t.tipo = tipo
	return t
}

func imprimirError(tok string, l, j int) {
	fmt.Print("Error Lexico: " + tok + " en la linea ")
	fmt.Print(l)
	fmt.Print(" columna ")
	fmt.Println(j)
}
