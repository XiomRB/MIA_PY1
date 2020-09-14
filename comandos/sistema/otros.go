package sistema

import (
	"Archivos/PY1/comandos/disco"
	"fmt"
	"io/ioutil"
	"strings"
)

type Ren struct {
	Id   string
	Name string
	Path string
}

type Cat struct {
	Id   string
	File string
}

func AdminRen(comando Ren) {
	if len(comando.Id) == 0 {
		fmt.Println("Error: el parametro id es obligatorio")
	} else if len(comando.Name) == 0 {
		fmt.Println("Error: el parametro name es obligatorio")
	} else if len(comando.Path) == 0 {
		fmt.Println("Error: el parametro path es obligatorio")
	} else {
		letra, indice := EncontrarMontada(comando.Id)
		if letra == -1 {
			fmt.Println("Error: la particion no ha sido montada")
			return
		}
		cambiarNombre(&disco.DiscosMontados[letra].Particiones[indice], comando)
	}
}
func cambiarNombre(part *disco.Montada, comando Ren) {
	comando.Name = ElimComillas(comando.Name)
	if strings.HasSuffix(comando.Path, ".txt") || strings.HasSuffix(comando.Path, ".txt\"") {
		fmt.Println(comando.Path)
		carpeta := SearchCarpeta(part, comando.Path)
		if carpeta == -1 {
			fmt.Println("Error: la ruta no existe")
			return
		}
		var name [20]byte
		copy(name[:], comando.Name)
		detalle, file := SearchArchivo(part, int(part.AVD[carpeta].IndiceDD), name)
		if file < 5 && detalle != -1 {
			fmt.Println("Error: Ya existe un archivo con ese nombre")
			return
		}
		p := ElimComillas(comando.Path)
		lista := strings.Split(p, "/")
		var anterior [20]byte
		copy(anterior[:], lista[len(lista)-1])
		fmt.Println(anterior)
		detalle, file = SearchArchivo(part, int(part.AVD[carpeta].IndiceDD), anterior)
		if detalle == -1 || file > 4 {
			fmt.Println("Error: el archivo no existe")
			return
		}
		copy(part.DD[detalle].Files[file].Nombre[:], comando.Name)
		fmt.Println("Nombre de archivo modificado")
	} else {
		comando.Path += "/" + comando.Name
		carpeta := SearchCarpeta(part, comando.Path)
		if carpeta == -1 {
			fmt.Println("Error: la ruta no existe")
			return
		}
		copy(part.AVD[carpeta].Nombre[:], comando.Name)
		fmt.Println("Nombre de carpeta modificado")
	}
}

func AdminCat(comando Cat) {
	if len(comando.Id) == 0 {
		fmt.Println("Error: el parametro id es obligatorio")
	} else if len(comando.File) == 0 {
		fmt.Println("Error: el parametro File1 es obligatorio")
	} else {
		letra, indice := EncontrarMontada(comando.Id)
		if letra == -1 {
			fmt.Println("Error: la particion no esta montada")
			return
		}
		reportarCat(&disco.DiscosMontados[letra].Particiones[indice], comando.File)
	}

}

func reportarCat(part *disco.Montada, arch string) {
	carpeta := SearchCarpeta(part, arch)
	if carpeta == -1 {
		fmt.Println("Error: ruta incorrecta")
		return
	}
	arch = ElimComillas(arch)
	var nombre [20]byte
	lista := strings.Split(arch, "/")
	copy(nombre[:], lista[len(lista)-1])
	detalle, file := SearchArchivo(part, int(part.AVD[carpeta].IndiceDD), nombre)
	if detalle == -1 || file > 4 {
		fmt.Println("Error: no existe el archivo")
		return
	}
	var inodos []int
	if part.DD[detalle].Files[file].Inodo != -1 {
		inodos = append(inodos, int(part.DD[detalle].Files[file].Inodo))
	}
	inodos = SearchInodos(part, inodos, int(part.DD[detalle].Files[file].Inodo))
	if len(inodos) == 0 {
		fmt.Println("El archivo no contiene texto")
		return
	}
	var contenido []byte
	var bloque int
	for i := 0; i < len(inodos); i++ {
		for j := 0; j < 4; j++ {
			bloque = int(part.Inodos[inodos[i]].Bloques[j])
			if bloque != -1 {
				for k := 0; k < 25; k++ {
					if part.BB[bloque].Text[k] != 0 {
						contenido = append(contenido, part.BB[bloque].Text[k])
					}
				}
			}
		}
	}
	ioutil.WriteFile("/home/gabriela/cat/"+lista[len(lista)-1], contenido, 0755)
	fmt.Println("Archivo creado en la carpeta cat")
}
