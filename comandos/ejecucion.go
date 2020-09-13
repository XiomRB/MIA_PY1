package comandos

import (
	"Archivos/PY1/analizador"
	"Archivos/PY1/comandos/disco"
	"Archivos/PY1/comandos/reportes"
	"Archivos/PY1/comandos/sistema"
	"Archivos/PY1/estructuras"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

//Ejecutar todo
func Ejecutar(cadena string) {
	if len(cadena) > 0 {
		raiz := analizador.Parser(cadena)
		for i := 0; i < len(raiz.Hijos); i++ {
			analizador.Imprimir(raiz.Hijos[i])
			leerComando(raiz.Hijos[i])
			fmt.Println("")
		}
	}
}

func leerComando(raiz analizador.Nodo) {
	switch strings.ToUpper(raiz.Dato) {
	case "MKDISK":
		mkdisk := disco.Mkdisk{-1, "", "", "m"}
		for i := 0; i < len(raiz.Hijos); i++ {
			validarMKDISK(raiz.Hijos[i], &mkdisk)
		}
		if estructuras.ValidarPath(mkdisk.Path, raiz.Linea) {
			if estructuras.VerificarName(mkdisk.Name, raiz.Linea) {
				if estructuras.VerificarSize(mkdisk.Size, raiz.Linea) {
					if len(mkdisk.Unit) > 0 {
						mkdisk.Size = estructuras.DarSize(mkdisk.Size, mkdisk.Unit)
						fmt.Println(disco.CrearDisco(mkdisk))
					}
				}
			}
		}
	case "RMDISK":
		n := raiz.Hijos[0]
		if strings.EqualFold(n.Dato, "path") {
			fmt.Println("Esta seguro que desea eliminar el disco ?\n  1.Si       2.No")
			elim := ""
			fmt.Scanln(&elim)
			if elim == "1" {
				disco.EliminarDisco(analizador.HomePath(n.Hijos[0]), raiz.Linea)
			}
		} else {
			fmt.Println("Error: El parametro path es obligatorio  --Linea: ", raiz.Linea)
		}
	case "FDISK":
		fdisk := disco.Fdisk{"", "", "w", "p", -1, "k", 0, ""}
		for i := 0; i < len(raiz.Hijos); i++ {
			validarFDISK(raiz.Hijos[i], &fdisk)
		}
		if estructuras.ValidarPath(fdisk.Path, raiz.Linea) {
			if estructuras.VerificarName(fdisk.Name, raiz.Linea) {
				if fdisk.Add == 0 && len(fdisk.Delete) > 1 { //Parametro delete
					if fdisk.Size == -1 {
						//---------------------------------------------------------------------------------------metodo que elimina
						disco.Administrar(fdisk)
					} else {
						fmt.Println("Error: los parametros size y delete no pueden ir juntos --Linea: ", raiz.Linea)
					}
				} else if (fdisk.Add != -1000000 && fdisk.Add != 0) && len(fdisk.Delete) == 0 { //parametro add
					if fdisk.Size == -1 {
						fdisk.Add = estructuras.DarSize(fdisk.Add, fdisk.Unit)
						disco.Administrar(fdisk)
						//--------------------------------------------------------------------------------------------metodo para aniade
					} else {
						fmt.Println("Error: los parametros size y add no pueden ir juntos --Linea: ", raiz.Linea)
					}
				} else if fdisk.Add == 0 && len(fdisk.Delete) == 0 {
					if estructuras.VerificarSize(fdisk.Size, raiz.Linea) {
						if len(fdisk.Unit) > 0 {
							fdisk.Size = estructuras.DarSize(fdisk.Size, fdisk.Unit)
							if len(fdisk.Tipo) > 0 {
								disco.Administrar(fdisk)
							}
							//---------------------------------------------------------------------------------------metodo para crear particion
						}
					}
				}
			}
		}
	case "MOUNT":
		mount := disco.Mount{}
		if len(raiz.Hijos) == 0 {
			disco.MostrarMontadas()
		} else {
			for i := 0; i < len(raiz.Hijos); i++ {
				validarMount(raiz.Hijos[i], &mount)
			}
			if len(mount.Path) == 0 || len(mount.Name) == 0 {
				fmt.Println("Error: Los parametros no corresponden al comando mount")
			} else {
				disco.Montar(mount)
			}
		}
	case "UNMOUNT":
		unmount := disco.Unmount{}
		for i := 0; i < len(raiz.Hijos); i++ {
			validarUnmount(raiz.Hijos[i], &unmount)
		}
		disco.Desmontar(unmount)
	case "MKFS":
		mkfs := sistema.Mkfs{"", "", 0, "k"}
		for i := 0; i < len(raiz.Hijos); i++ {
			validarMKFS(raiz.Hijos[i], &mkfs)
		}
		sistema.AdminComando(mkfs)
	case "LOGIN":
		login := sistema.Login{}
		for i := 0; i < len(raiz.Hijos); i++ {
			validarLoguin(raiz.Hijos[i], &login)
		}
		sistema.AdminLogin(login)
	case "LOGOUT":
		if len(raiz.Hijos) != 0 {
			fmt.Println("Error: Este comando no recibe parametros")
		} else {
			if sistema.LoginUs.Estado {
				sistema.LoginUs.Estado = false
			} else {
				fmt.Println("Error: No hay un usuario logueado")
			}

		}
	case "MKGRP":
		grp := sistema.Grp{}
		grp.Accion = true
		for i := 0; i < len(raiz.Hijos); i++ {
			validargrup(raiz.Hijos[i], &grp)
		}
		sistema.AdminGrupos(grp)
	case "RMGRP":
		grp := sistema.Grp{}
		grp.Accion = false
		for i := 0; i < len(raiz.Hijos); i++ {
			validargrup(raiz.Hijos[i], &grp)
		}
		sistema.AdminGrupos(grp)
	case "MKUSR":
		mkusr := sistema.Mkusr{}
		for i := 0; i < len(raiz.Hijos); i++ {
			validarmkusr(raiz.Hijos[i], &mkusr)
		}
		fmt.Println("entro en mjusr ", mkusr.Usr)
		sistema.CrearUsuario(mkusr)
	case "RMUSR":
		rmusr := sistema.Rmusr{}
		for i := 0; i < len(raiz.Hijos); i++ {
			validarrmusr(raiz.Hijos[i], &rmusr)
		}
		sistema.EliminarUsuario(rmusr)
	case "CHMOD":
	case "MKFILE":
		mkfile := sistema.Mkfile{}
		for i := 0; i < len(raiz.Hijos); i++ {
			validarMkfile(raiz.Hijos[i], &mkfile)
		}
		sistema.AdminMkFile(mkfile)
	case "CAT":
	case "RM":
	case "EDIT":
	case "REN":
	case "MKDIR":
		mkdir := sistema.Mkdir{}
		for i := 0; i < len(raiz.Hijos); i++ {
			validarmkdir(raiz.Hijos[i], &mkdir)
		}
		sistema.AdminCarpetas(mkdir)
	case "CP":
	case "MV":
	case "FIND":
	case "CHOWN":
	case "CHGRP":
	case "LOSS":
	case "RECOVERY":
	case "REP":
		reporte := reportes.Reporte{}
		for i := 0; i < len(raiz.Hijos); i++ {
			validarRep(raiz.Hijos[i], &reporte)
		}
		reportes.AdministrarReportes(reporte)
	case "PAUSE":
		var pausa string
		fmt.Println("El programa esta en pausa, presione cualquier letra para continuar")
		fmt.Scanln(&pausa)
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
	default:
		fmt.Println("Error: comando no reconocido")
	}
}

func validarMKDISK(raiz analizador.Nodo, comando *disco.Mkdisk) {
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
		comando.Unit = strings.ToLower(analizador.ValidarUnidad(false, raiz.Hijos[0]))
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
			log.Fatal(err)
		} else {
			comando.Add = int64(s)
		}
	case "delete":
		if strings.EqualFold(raiz.Hijos[0].Dato, "fast") || strings.EqualFold(raiz.Hijos[0].Dato, "full") {
			comando.Delete = strings.ToLower(raiz.Hijos[0].Dato)
		} else {
			comando.Delete = "e"
		}
	case "unit":
		comando.Unit = analizador.ValidarUnidad(true, raiz.Hijos[0])
	case "type":
		tipo := raiz.Hijos[0].Dato
		if strings.EqualFold(tipo, "p") || strings.EqualFold(tipo, "e") || strings.EqualFold(tipo, "l") {
			comando.Tipo = strings.ToLower(tipo)
		} else {
			comando.Tipo = ""
			fmt.Println("Error: el parametro tipo solo acepta p, l o e")
		}
	case "fit":
		comando.Fit = strings.ToLower(analizador.ValidarFit(raiz.Hijos[0]))
	}
}

func validarMount(raiz analizador.Nodo, comando *disco.Mount) {
	switch strings.ToLower(raiz.Dato) {
	case "path":
		comando.Path = analizador.HomePath(raiz.Hijos[0])
	case "name":
		comando.Name = raiz.Hijos[0].Dato
	}
}

func validarUnmount(raiz analizador.Nodo, comando *disco.Unmount) {
	if strings.EqualFold(raiz.Dato, "id") {
		comando.Desmontadas = append(comando.Desmontadas, raiz.Hijos[0].Dato)
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

func validarRep(raiz analizador.Nodo, comando *reportes.Reporte) {
	switch strings.ToLower(raiz.Dato) {
	case "path":
		comando.Path = analizador.HomePath(raiz.Hijos[0])
	case "name":
		comando.Name = raiz.Hijos[0].Dato
	case "id":
		comando.Id = raiz.Hijos[0].Dato
	case "ruta":
		comando.Ruta = raiz.Hijos[0].Dato
	}
}

func validarMKFS(raiz analizador.Nodo, comando *sistema.Mkfs) {
	switch strings.ToLower(raiz.Dato) {
	case "id":
		comando.Id = raiz.Hijos[0].Dato
	case "add":
		s, err := strconv.Atoi(raiz.Hijos[0].Dato)
		if err != nil {
			log.Fatal(err)
		}
		comando.Add = int64(s)
	case "type":
		comando.Tipo = raiz.Hijos[0].Dato
	case "unit":
		comando.Unit = analizador.ValidarUnidad(true, raiz.Hijos[0])

	}
}

func validarLoguin(raiz analizador.Nodo, loguin *sistema.Login) {
	switch strings.ToLower(raiz.Dato) {
	case "id":
		loguin.Id = raiz.Hijos[0].Dato
	case "usr":
		loguin.Usr = raiz.Hijos[0].Dato
	case "pwd":
		loguin.Pwd = raiz.Hijos[0].Dato
	}
}

func validargrup(raiz analizador.Nodo, comando *sistema.Grp) {
	switch strings.ToLower(raiz.Dato) {
	case "id":
		comando.Id = raiz.Hijos[0].Dato
	case "name":
		comando.Name = raiz.Hijos[0].Dato
	}
}

func validarmkusr(raiz analizador.Nodo, comando *sistema.Mkusr) {
	switch strings.ToLower(raiz.Dato) {
	case "id":
		comando.Id = raiz.Hijos[0].Dato
	case "grp":
		comando.Grp = raiz.Hijos[0].Dato
	case "pwd":
		comando.Pwd = raiz.Hijos[0].Dato
	case "usr":
		comando.Usr = raiz.Hijos[0].Dato
	}
}

func validarrmusr(raiz analizador.Nodo, comando *sistema.Rmusr) {
	switch strings.ToLower(raiz.Dato) {
	case "id":
		comando.Id = raiz.Hijos[0].Dato
	case "usr":
		comando.Usr = raiz.Hijos[0].Dato
	}
}

func validarmkdir(raiz analizador.Nodo, mkdir *sistema.Mkdir) {
	switch strings.ToLower(raiz.Dato) {
	case "id":
		mkdir.Id = raiz.Hijos[0].Dato
	case "path":
		mkdir.Path = raiz.Hijos[0].Dato
	case "p":
		mkdir.Padre = true
	}
}

func validarchmod(raiz analizador.Nodo, chmod *sistema.Chmod) {
	switch strings.ToLower(raiz.Dato) {
	case "id":
		chmod.Id = raiz.Hijos[0].Dato
	case "path":
		chmod.Path = raiz.Hijos[0].Dato
	case "r":
		chmod.R = true
	case "ugo":
		chmod.Ugo = raiz.Hijos[0].Dato
	}
}

func validarMkfile(raiz analizador.Nodo, mkfile *sistema.Mkfile) {
	switch strings.ToLower(raiz.Dato) {
	case "id":
		mkfile.Id = raiz.Hijos[0].Dato
	case "path":
		mkfile.Path = raiz.Hijos[0].Dato
	case "p":
		mkfile.P = true
	case "size":
		s, err := strconv.Atoi(raiz.Hijos[0].Dato)
		if err != nil {
			mkfile.Size = 0
			return
		}
		mkfile.Size = int64(s)
	case "cont":
		mkfile.Cont = raiz.Hijos[0].Dato
	}
}
