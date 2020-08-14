package comandos

import (
	"Archivos/PY1/analizador"
	"fmt"
	"strconv"
	"strings"
)

func ejecutarCom(cadena string) {
	raiz := analizador.Parser(cadena)
}

func leerComando(raiz analizador.Nodo) {
	switch strings.ToUpper(raiz.Dato) {
	case "MKDISK":
		var mkdisk Mkdisk
		for i := 0; i < len(raiz.Hijos); i++ {
			validarMKDISK(raiz.Hijos[i], mkdisk)
		}

	case "RMDISK":
	case "FDISK":
	case "MOUNT":
	case "UNMOUNT":
	case "MKFS":
	case "LOGIN":
	case "LOGOUT":
	case "MKGRP":
	case "RMGRP":
	case "MKUSR":
	case "RMUSR":
	case "CHMOD":
	case "MKFILE":
	case "CAT":
	case "RM":
	case "EDIT":
	case "REN":
	case "MKDIR":
	case "CP":
	case "MV":
	case "FIND":
	case "CHOWN":
	case "CHGRP":
	case "LOSS":
	case "RECOVERY":
	case "REP":
	case "PAUSE":
	case "EXEC":
	}
}

func validarMKDISK(raiz analizador.Nodo, comando Mkdisk) {

}

func nuevoPath(p analizador.Nodo) string {
	path := ""
	if p.Dato[0] == 34 {
		for i := 1; i < len(p.Dato); i++ {
			if p.Dato[i] == 32 {
				path += "_"
			} else {
				path += string(p.Dato[i])
			}
		}
	}
	return p.Dato
}

func validarSize(f analizador.Nodo) int {
	s, err := strconv.Atoi(f.Dato)
	if err == nil {
		if s > 0 {
			return s
		}
		fmt.Println("Error: El parametro size debe ser mayor a 0  ---Linea: %l Col: %c", f.Linea, f.Columna)
		return 0
	}
	fmt.Println("Error: El parametro size solo reconoce numeros  ---Linea: %l Col: %c", f.Linea, f.Columna)
	return 0
}

func validarFit(f analizador.Nodo) string {
	if strings.EqualFold(f.Dato, "FF") || strings.EqualFold(f.Dato, "WF") || strings.EqualFold(f.Dato, "BF") {
		return string(f.Dato[0])
	}
	fmt.Println("Error: El parametro fit no reconoce ese valor ---Linea: %l Col: %c", f.Linea, f.Columna)
	return ""
}

func verificarUnidad(part bool, u analizador.Nodo) string {
	if strings.EqualFold(u.Dato, "k") || strings.EqualFold(u.Dato, "m") {
		return string(u.Dato[0])
	}
	if strings.EqualFold(u.Dato, "b") && part {
		return string(u.Dato[0])
	}
	fmt.Println("Error: El parametro unit no reconoce ese valor ---Linea: %l Col: %c", u.Linea, u.Columna)
	return ""
}

func darSize(s int, u rune) int {
	if strings.EqualFold(string(u), "k") {
		return s * 1024
	} else if strings.EqualFold(string(u), "m") {
		return s * 1024 * 1024
	}
	return s
}

/*func validarNombre(n analizador.Nodo) bool{
	if strings.HasSuffix(n.Dato,".dsk"){
		return true
	}
	return false
}*/
