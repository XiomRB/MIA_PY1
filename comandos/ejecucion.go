package comandos

import (
	"Archivos/PY1/analizador"
	"Archivos/PY1/comandos/disco"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

//Ejecutar todo
func Ejecutar(cadena string) {
	if len(cadena) > 0 {
		raiz := analizador.Parser(cadena)
		for i := 0; i < len(raiz.Hijos); i++ {
			leerComando(raiz.Hijos[i])
		}
	}
}

func leerComando(raiz analizador.Nodo) {
	switch strings.ToUpper(raiz.Dato) {
	case "MKDISK":
		mkdisk := Mkdisk{-1, "", "", "m"}
		for i := 0; i < len(raiz.Hijos); i++ {
			validarMKDISK(raiz.Hijos[i], &mkdisk)
		}
		if ValidarPath(mkdisk.Path, raiz.Linea) {
			if VerificarName(mkdisk.Name, raiz.Linea) {
				if VerificarSize(mkdisk.Size, raiz.Linea) {
					if len(mkdisk.Unit) > 0 {
						mkdisk.Size = DarSize(mkdisk.Size, mkdisk.Unit)
						fmt.Println(CrearDisco(mkdisk))
					}
				}
			}
		}
	case "RMDISK":
		n := raiz.Hijos[0]
		if strings.EqualFold(n.Dato, "path") {
			eliminarDisco(analizador.HomePath(n.Hijos[0]), raiz.Linea)
		} else {
			fmt.Println("Error: El parametro path es obligatorio  --Linea: ", raiz.Linea)
		}
	case "FDISK":
		fdisk := disco.Fdisk{"", "", "w", "p", -1, "k", 0, ""}
		for i := 0; i < len(raiz.Hijos); i++ {
			validarFDISK(raiz.Hijos[i], &fdisk)
		}
		if ValidarPath(fdisk.Path, raiz.Linea) {
			if VerificarName(fdisk.Name, raiz.Linea) {
				if fdisk.Add == 0 && len(fdisk.Delete) > 1 { //Parametro delete
					if fdisk.Size == 0 {
						if len(fdisk.Unit) > 0 {
							fdisk.Add = DarSize(fdisk.Add, fdisk.Unit)
							//---------------------------------------------------------------------------------------metodo que aniade
						}
					} else {
						fmt.Println("Error: los parametros size y add no pueden ir juntos --Linea: ", raiz.Linea)
					}
				} else if (fdisk.Add != -1000000 && fdisk.Add != 0) && len(fdisk.Delete) == 0 { //parametro add
					if fdisk.Size == 0 {
						//--------------------------------------------------------------------------------------------metodo para eliminar
					} else {
						fmt.Println("Error: los parametros size y delete no pueden ir juntos --Linea: ", raiz.Linea)
					}
				} else if fdisk.Add == 0 && len(fdisk.Delete) == 0 {
					if VerificarSize(fdisk.Size, raiz.Linea) {
						if len(fdisk.Unit) > 0 {
							fdisk.Size = DarSize(fdisk.Size, fdisk.Unit)
							//---------------------------------------------------------------------------------------metodo para crear particion
						}
					}
				}
			}
		}
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
		fmt.Println("El programa esta en pausa, presione cualquier letra para continuar")
		var tecla string
		for len(tecla) == 0 {
			fmt.Scan(&tecla)
		}
	case "EXEC":
		p := raiz.Hijos[0]
		if strings.EqualFold(p.Dato, "path") {
			path := analizador.HomePath(p.Hijos[0])
			if strings.HasSuffix(path, ".mia") {
				archivo := ejecutarArchivo(path)
				if len(archivo) > 0 {
					Ejecutar(archivo)
				}
			} else {
				fmt.Println("Error: El archivo que desea abrir no es de extension mia")
			}
		}
	}
}

func validarMKDISK(raiz analizador.Nodo, comando *Mkdisk) {
	switch strings.ToLower(raiz.Dato) {
	case "path":
		comando.Path = analizador.HomePath(raiz.Hijos[0])
	case "name":
		if strings.HasSuffix(raiz.Hijos[0].Dato, ".dsk") {
			comando.Name = raiz.Hijos[0].Dato
		} else {
			comando.Name = "0"
			fmt.Println("Error: El parametro name no posee la extension .disk   --- Linea: ", raiz.Linea, " Col: ", raiz.Columna)
		}
	case "size":
		comando.Size = analizador.ValidarSize(raiz.Hijos[0])
	case "unit":
		comando.Unit = analizador.ValidarUnidad(false, raiz.Hijos[0])
	}
}

func validarFDISK(raiz analizador.Nodo, comando *disco.Fdisk) {
	switch strings.ToLower(raiz.Dato) {
	case "path":
		comando.Path = analizador.HomePath(raiz.Hijos[0])
	case "name":
		comando.Name = raiz.Hijos[0].Dato
	case "size":
		comando.Size = analizador.ValidarSize(raiz.Hijos[0])
	case "add":
		s, err := strconv.Atoi(raiz.Hijos[0].Dato)
		if err != nil {
			fmt.Println("Error: el parametro add solo recibe valores numericos --Linea: ", raiz.Linea)
			comando.Add = -1000000
		} else {
			comando.Add = s
		}
	case "delete":
		if strings.EqualFold(raiz.Hijos[0].Dato, "fast") || strings.EqualFold(raiz.Hijos[0].Dato, "full") {
			comando.Delete = raiz.Hijos[0].Dato
		} else {
			comando.Delete = "e"
		}
	case "unit":
		comando.Unit = analizador.ValidarUnidad(true, raiz.Hijos[0])
	case "type":
		tipo := raiz.Hijos[0].Dato
		if strings.EqualFold(tipo, "p") || strings.EqualFold(tipo, "e") || strings.EqualFold(tipo, "l") {
			comando.Tipo = tipo
		} else {
			comando.Tipo = "f"
		}
	case "fit":
		comando.Fit = analizador.ValidarFit(raiz.Hijos[0])
	}
}

func ejecutarArchivo(n string) string {
	data, err := ioutil.ReadFile(n)
	if err != nil {
		fmt.Println("Hubo un error al leer el archivo")
		return ""
	}
	return string(data)
}
