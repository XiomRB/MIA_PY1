package sistema

import (
	"Archivos/PY1/comandos/disco"
	"Archivos/PY1/estructuras"
	"fmt"
)

var LoginUs estructuras.Logueado

type Mkusr struct {
	Id  string
	Usr string
	Pwd string
	Grp string
}

type Rmusr struct {
	Id  string
	Usr string
}

type Login struct {
	Usr string
	Pwd string
	Id  string
}

type Grp struct {
	Accion bool
	Id     string
	Name   string
}

func crearGrupo(name string, particion *disco.Montada) {
	var n [10]byte
	copy(n[:], name)
	for i := 0; i < len(particion.Grupos); i++ {
		if particion.Grupos[i].Name == n && particion.Grupos[i].Estado {
			fmt.Println("Error: el grupo ya existe")
			return
		}
	}
	grupo := estructuras.Grupo{}
	grupo.Indice = len(particion.Grupos) + 1
	grupo.Name = n
	grupo.Estado = true
	particion.Grupos = append(particion.Grupos, grupo)
	fmt.Println("Grupo creado")
}

func eliminarGrupo(name string, particion *disco.Montada) {
	var n [10]byte
	copy(n[:], name)
	for i := 0; i < len(particion.Grupos); i++ {
		if particion.Grupos[i].Name == n && particion.Grupos[i].Estado {
			particion.Grupos[i].Estado = false
			fmt.Println("Grupo eliminado")
			return
		}
	}
	fmt.Println("Error: el grupo no existe")
}

func EliminarUsuario(comando Rmusr) {
	var root [10]byte
	copy(root[:], "root")
	var usr [10]byte
	copy(usr[:], comando.Usr)
	if LoginUs.Estado && LoginUs.Name == root {
		letra, indice, path := EncontrarMontada(comando.Id)
		if len(path) == 0 {
			fmt.Println("Error, la particion no ha sido montada")
		} else {
			for i := 0; i < len(disco.DiscosMontados[letra].Particiones[indice].Grupos); i++ {
				if disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Estado {
					for j := 0; j < len(disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Usuarios); j++ {
						if disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Usuarios[j].Name == usr && disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Usuarios[j].Estado {
							disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Usuarios[j].Estado = false
							fmt.Println("Usuario eliminado")
							return
						}
					}
				}
			}
			fmt.Println("Error: el usuario no existe")
		}
	} else {
		fmt.Println("Error: el comando solo puede ser usando por un usuario root")
	}
}

func CrearUsuario(comando Mkusr) {
	var root [10]byte
	usuario := estructuras.Usuario{}
	usuario.Estado = true
	copy(root[:], "root")
	var grp [10]byte
	copy(grp[:], comando.Grp)
	copy(usuario.Name[:], comando.Usr)
	copy(usuario.Clave[:], comando.Pwd)
	if LoginUs.Estado && LoginUs.Name == root {
		letra, indice, path := EncontrarMontada(comando.Id)
		if len(path) == 0 {
			fmt.Println("Error, la particion no ha sido montada")
		} else {
			for i := 0; i < len(disco.DiscosMontados[letra].Particiones[indice].Grupos); i++ {
				if disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Name == grp && disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Estado {
					for j := 0; j < len(disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Usuarios); j++ {
						if disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Usuarios[j].Name == usuario.Name && disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Usuarios[j].Estado {
							fmt.Println("Error: el usuario ya existe")
							return
						}
					}
					disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Usuarios = append(disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Usuarios, usuario)
					fmt.Println("Usuario creado")
					return
				}
			}
			fmt.Println("Error: el grupo no existe, no se puede crear el usuario")
		}
	} else {
		fmt.Println("Error: Debe estar logueado con el usuario root para poder usar este comando")
	}
}

func loguearse(comando Login) {
	if LoginUs.Estado {
		fmt.Println("Error: Hay un usuario que aun no ha cerrado sesion")
	} else {
		letra, indice, path := EncontrarMontada(comando.Id)
		if len(path) == 0 {
			fmt.Println("Error: La particion no se encuentra montada")
		} else {
			var usr [10]byte
			copy(usr[:], comando.Usr)
			for i := 0; i < len(disco.DiscosMontados[letra].Particiones[indice].Grupos); i++ {
				for j := 0; j < len(disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Usuarios); j++ {
					if usr == disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Usuarios[j].Name && disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Estado && disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Usuarios[j].Estado {
						var pass [10]byte
						copy(pass[:], comando.Pwd)
						if pass == disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Usuarios[j].Clave {
							LoginUs.Name = usr
							LoginUs.Estado = true
							LoginUs.Particion = comando.Id
							LoginUs.Grupo = disco.DiscosMontados[letra].Particiones[indice].Grupos[i].Name
							fmt.Println("Usuario Logueado")
							return
						}
					}
				}
			}
			fmt.Println("Error: El usuario no existe, revise que introdujo bien sus datos")
		}
	}
}

func AdminLogin(login Login) {
	if len(login.Id) == 0 {
		fmt.Println("Error: El parametro id es obligatorio")
	} else if len(login.Usr) == 0 {
		fmt.Println("Error: El parametro usr es obligatorio")
	} else if len(login.Pwd) == 0 {
		fmt.Println("Error: El parametro pwd es obligatorio")
	} else {
		loguearse(login)
	}
}

func AdminGrupos(comando Grp) {
	var usr [10]byte
	copy(usr[:], "root")
	if usr == LoginUs.Name && LoginUs.Estado {
		if len(comando.Id) != 0 {
			if len(comando.Name) != 0 {
				letra, indice, path := EncontrarMontada(comando.Id)
				if len(path) == 0 {
					fmt.Println("Error: La particion no ha sido montada")
				} else {
					if comando.Accion { //si el comando es mkgrp
						crearGrupo(comando.Name, &disco.DiscosMontados[letra].Particiones[indice])
					} else { //si el comando es rmgrp
						eliminarGrupo(comando.Name, &disco.DiscosMontados[letra].Particiones[indice])
					}
				}
			} else {
				fmt.Println("Error: el parametro name es obligatorio")
			}
		} else {
			fmt.Println("Error: el parametro id es obligatorio")
		}
	} else {
		fmt.Println("Error: El usuario root es el unico que puede utilizar este comando")
	}
}
