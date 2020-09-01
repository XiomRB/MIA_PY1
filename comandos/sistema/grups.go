package sistema

import (
	"Archivos/PY1/estructuras"
	"fmt"
)

var LoginUs estructuras.Logueado

type Login struct {
	Usr string
	Pwd string
	Id  string
}

func crearNuevoGrupo() {

}

func crearNuevoUsuario() {

}

func loguearse(comando Login) {
	if LoginUs.Estado {
		fmt.Println("Error: Hay un usuario que aun no ha cerrado sesion")
	} else {
		particion, path := EncontrarMontada(comando.Id)
		if len(path) == 0 {
			fmt.Println("Error: La particion no se encuentra montada")
		} else {
			var usr [10]byte
			copy(usr[:], comando.Usr)
			for i := 0; i < len(particion.Grupos); i++ {
				for j := 0; j < len(particion.Grupos[i].Usuarios); j++ {
					if usr == particion.Grupos[i].Usuarios[j].Name {
						var pass [10]byte
						copy(pass[:], comando.Pwd)
						if pass == particion.Grupos[i].Usuarios[j].Clave {
							LoginUs.Name = usr
							LoginUs.Estado = true
							LoginUs.Particion = comando.Id
							LoginUs.Grupo = particion.Grupos[i].Name
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
