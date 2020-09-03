package sistema

import "Archivos/PY1/estructuras"

type Chmod struct {
	Id   string
	Path string
	Ugo  int
	R    [3]byte
}

func AdminChmod(comadno Chmod) {

}

func VerificarPermisos(permisos [3]byte, prop estructuras.Propietario) {

}
