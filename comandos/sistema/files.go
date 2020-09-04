package sistema

import (
	"Archivos/PY1/comandos/disco"
	"Archivos/PY1/estructuras"
	"fmt"
	"strings"
)

type Mkfile struct {
	Id   string
	Path string
	P    bool
	Size int64
	Cont string
}

func AdminMkFile(comando Mkfile) {
	if len(comando.Id) == 0 {
		fmt.Println("Error: el parametro id es obligatorio")
	} else if len(comando.Path) == 0 {
		fmt.Println("Error: el parametro path es obligatorio")
	} else {
		letra, indice, path := EncontrarMontada(comando.Id)
	}
}

func CrearArchivo(path string, part *disco.Montada, crear bool) {
	direc := ElimComillas(path)
	lista := strings.Split(direc, "/")
	direc = "/"
	for i := 1; i < len(lista)-1; i++ {
		direc += lista[i]
	}
	var nameFile [20]byte
	copy(nameFile[:], lista[len(lista)-1])
	carpeta := BuscarCarpeta(direc, part, crear)

}

func EscribirArchivo() {

}

func EncontrarDetalle(arch [20]byte, part *disco.Montada, indice int) (int, int) { //retorna detalle e indice de archivo en el detalle
	var n [20]byte
	i := 0
	for i = 0; i < len(part.DD[indice].Files); i++ {
		if part.DD[indice].Files[i].Nombre == arch || part.DD[indice].Files[i].Nombre == n {
			return indice, i
		}
	}
	if part.DD[indice].Next > 0 {
		indice, i = EncontrarDetalle(arch, part, int(part.DD[indice].Next))
	} else {
		nuevoDet := estructuras.DetalleDir{}
		j := 0
		for j = 0; j < len(part.BitmapDetalle); j++ {
			if part.BitmapDetalle[j] == 0 {
				part.DD[j] = nuevoDet
				part.BitmapDetalle[j] = 1
				break
			}
		}
		if j == len(part.BitmapDetalle) {
			part.DD = append(part.DD, nuevoDet)
			part.BitmapDetalle = append(part.BitmapDetalle, 1)
		}
		part.DD[indice].Next = int64(j)
		i = 0
		indice = j
	}
	return indice, i
}
