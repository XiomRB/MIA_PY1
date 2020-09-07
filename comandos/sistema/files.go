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

func CrearArchivo(comando Mkfile, part *disco.Montada, crear bool) {
	direc := ElimComillas(comando.Path)
	lista := strings.Split(direc, "/")
	direc = "/"
	for i := 1; i < len(lista)-1; i++ {
		direc += lista[i]
	}
	var nameFile [20]byte
	copy(nameFile[:], lista[len(lista)-1])
	carpeta := BuscarCarpeta(direc, part, crear) //se busca la carpeta donde se creara el archivo
	indiceDet := 0
	indiceFile := 0
	if carpeta >= 0 {
		if part.AVD[carpeta].IndiceDD != -1 {
			indiceDet, indiceFile = EncontrarDetalle(nameFile, part, int(part.AVD[carpeta].IndiceDD)) //encuentra el detalle de directorio de la carpeta elegida
		} else {
			nuevoDet := estructuras.DetalleDir{}
			nuevoDet.Next = -1
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
			part.AVD[carpeta].IndiceDD = int64(j)
			indiceDet = j
			indiceFile = 0
		}
		if part.DD[indiceDet].Files[indiceFile].Nombre == nameFile {
			fmt.Println("El archivo ya existe")
			return
		}
		part.DD[indiceDet].Files[indiceFile] = CrearFile(lista[len(lista)-1], -1)
		nb := 0
		if comando.Size > 0 {
			if comando.Size%100 == 0 { //se divide entre 100 ya que son 100 caracteres los que le caben a cada inodo
				nb = int(comando.Size) / 100
			} else {
				nb = (int(comando.Size) / 100) + 1
			}
		} else if len(comando.Cont) > 0 {
			if len(comando.Cont)%100 == 0 {
				nb = len(comando.Cont) / 100
			} else {
				nb = (len(comando.Cont) / 100) + 1
			}
		}
		if nb > 0 { //se crean los inodos necesarios para el archivo y se enlazan
			ids := make([]int, nb)
			if comando.Size > 0 {
				for i := 0; i < nb; i++ {
					ids[i] = CrearNInodo(part, comando.Size)
				}
			} else {
				for i := 0; i < nb; i++ {
					ids[i] = CrearNInodo(part, int64(len(comando.Cont)))
				}
			}
			for i := 0; i < nb-1; i++ {
				part.Inodos[i].Indirecto = part.Inodos[i+1].Indice
			}
			longitud := int(comando.Size) - len(comando.Cont)
			abc := "abcdefghijklmnopqrstuvwxyz"
			l := 0
			for i := len(comando.Cont); i < longitud; i++ { //se crea la cadena a escribir en el archivo
				if l == len(abc)-1 {
					l = 0
				}
				comando.Cont += string(abc[l])
				l++
			}
			EscribirArchivo(comando, part.Inodos[ids[0]].NBloques, part, ids)
			part.DD[indiceDet].Files[indiceFile].Inodo = int64(ids[0])
		}
		fmt.Println("Archivo creado")
	}
}

func EscribirArchivo(comando Mkfile, nb int64, part *disco.Montada, inodos []int) {
	bloques := EscribirBloques(comando.Cont, int64(nb))
	b := 0
	bb := make([]int, len(bloques))
	for i := 0; i < len(bloques); i++ {
		for b < len(part.BitmapBloques) {
			if part.BitmapBloques[b] == 0 {
				part.BitmapBloques[b] = 1
				part.BB[b] = bloques[i]
				bb[i] = b
				break
			}
			b++
		}
		if b >= len(part.BitmapBloques) {
			part.BitmapBloques = append(part.BitmapBloques, 1)
			part.BB = append(part.BB, bloques[i])
			bb[i] = len(part.BitmapBloques) - 1
			b++
		}
	}
	block := 0
	ids := 0
	for i := 0; i < len(bb); i++ {
		if i > 0 && i%4 == 0 {
			block = 0
			ids++
		}
		part.Inodos[inodos[ids]].Bloques[block] = int64(bb[i])
		block++
	}
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
		nuevoDet.Next = -1
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

func CrearNInodo(part *disco.Montada, size int64) int {
	inodo := CrearInodo(0, size)
	i := 0
	for i = 0; i < len(part.BitmapInodo); i++ {
		if part.BitmapInodo[i] == 0 {
			part.BitmapInodo[i] = 1
			part.Inodos[i] = inodo
			break
		}
	}
	if i == len(part.BitmapInodo) {
		part.BitmapInodo = append(part.BitmapInodo, 1)
		part.Inodos = append(part.Inodos, inodo)
	}
	part.Inodos[i].Indice = int64(i)
	return i
}
