package sistema

type Chmod struct {
	Id   string
	Path string
	Ugo  string
	R    bool
}

func AdminChmod(comadno Chmod) {

}

func VerificarPermisos(permiso byte) string {
	switch permiso {
	case 0:
		return "---"
	case 1:
		return "--x"
	case 2:
		return "-w-"
	case 3:
		return "-wx"
	case 4:
		return "r--"
	case 5:
		return "r-x"
	case 6:
		return "rw-"
	case 7:
		return "rwx"
	}
	return ""
}
