package sistema

import (
	"Archivos/PY1/comandos/disco"
	"fmt"
	"strings"
)

type Mkdir struct {
	Id    string
	Path  string
	Padre bool
}

func CrearCarpeta(particion *disco.Montada, nombre string, padre int) int { //el return es para cuando se crea un file
	index := 0
	apa := 0
	if LoginUs.Name == particion.AVD[padre].Prop.Name { //verifica que sea el propietario
		index, apa = EncontrarSubVacio(padre, particion)
		if apa == -1 {
			return apa
		}
		nvd := NuevoDirectorio(particion, nombre, apa)
		if nvd == -1 {
			return -1
		}
		particion.AVD[padre].IndicesSubs[index] = int64(nvd)
		return nvd
	} else if LoginUs.Grupo == particion.AVD[padre].Prop.Grupo {
		if strings.Contains(VerificarPermisos(particion.AVD[padre].Permisos[1]), "w") {
			index, apa = EncontrarSubVacio(padre, particion)
			if apa == -1 {
				return apa
			}
			nvd := NuevoDirectorio(particion, nombre, apa)
			if nvd == -1 {
				return -1
			}
			particion.AVD[padre].IndicesSubs[index] = int64(nvd)
			return nvd
		} else {
			fmt.Println("Error: no tiene los permisos, para crear carpetas, en la ruta especificada")
			return -1
		}
	} else if strings.Contains(VerificarPermisos(particion.AVD[padre].Permisos[2]), "w") {
		index, apa = EncontrarSubVacio(padre, particion)
		if apa == -1 {
			return apa
		}
		nvd := NuevoDirectorio(particion, nombre, apa)
		if nvd == -1 {
			return -1
		}
		particion.AVD[padre].IndicesSubs[index] = int64(nvd)
		return nvd
	} else {
		fmt.Println("Error: no tiene los permisos, para crear carpetas, en la ruta especificada")
		return -1
	}
}

func NuevoDirectorio(part *disco.Montada, nombre string, padre int) int {
	avd := CrearAVD(nombre)
	i := 0
	for i = 0; i < len(part.BitmapAVD); i++ {
		if part.BitmapAVD[i] == 0 {
			if i == len(part.AVD) {
				part.AVD = append(part.AVD, avd)
			} else {
				part.AVD[i] = avd //si hay uno directorio vacio se introduce el nuevo
			}
			part.BitmapAVD[i] = 1
			break
		}
	}
	if i == len(part.BitmapAVD) {
		fmt.Println("Error, ya creo el numero maximo de directorios")
		return -1
	}
	return i
}

func AdminCarpetas(comando Mkdir) {
	if len(comando.Id) == 0 {
		fmt.Println("Error: el parametro id es obligatorio")
	} else if len(comando.Path) == 0 {
		fmt.Println("Error: el parametro path es obligatorio")
	} else {
		if LoginUs.Estado {
			letra, indice := EncontrarMontada(comando.Id)
			if letra != -1 {
				decision := BuscarCarpeta(comando.Path, &disco.DiscosMontados[letra].Particiones[indice], comando.Padre)
				if decision == -2 {
					fmt.Println("Error: Hay carpetas padre que no existe")
				}
			} else {
				fmt.Println("Error: la particion no esta montada")
			}
		} else {
			fmt.Println("Error: debe estar logueado para poder realizar la accion")
		}
	}
}

func BuscarCarpeta(path string, particion *disco.Montada, crear bool) int {
	p := ElimComillas(path)
	lista := DescomponerRuta(p)
	i := 0
	indice := 0
	padre := 0
	for i = 1; i < len(lista)-1; i++ {
		indice, padre = EncontrarCarpeta(particion, lista[i], indice)
		if indice == 0 { //si no se encontro la carpeta se sale del for
			break
		}
	}
	if i < len(lista)-1 && crear { //al no encontrar carpeta
		for j := i; j < len(lista); j++ { //por cada nombre que no encontro se crea una nueva carpeta
			padre = CrearCarpeta(particion, lista[j], padre) //retorna el indice de la carpeta creada, para agregarle los hijos
			if padre == -1 {                                 //por si no tiene permiso de crear carpetas
				break
			}
		}
	} else if i == len(lista)-1 {
		padre = CrearCarpeta(particion, lista[len(lista)-1], indice)
	} else {
		return -2
	}
	return padre
}

func DescomponerRuta(path string) []string {
	lista := strings.Split(path, "/")
	lista[0] = "/"
	return lista
}

func ElimComillas(p string) string {
	path := ""
	if p[0] == 34 {
		for i := 1; i < len(p)-1; i++ {
			path += string(p[i])
		}
		return path
	}
	return p
}

func EncontrarSubVacio(padre int, part *disco.Montada) (int, int) {
	i := 0
	for i = 0; i < 6; i++ {
		if part.AVD[padre].IndicesSubs[i] == 0 {
			return i, padre
		}
	}
	if i == 6 {
		if part.AVD[padre].IndiceNext > 0 {
			i, padre = EncontrarSubVacio(int(part.AVD[padre].IndiceNext), part)
		} else {
			avd := CrearAVD(string(part.AVD[padre].Nombre[:]))
			avd.Permisos = part.AVD[padre].Permisos
			avd.Prop = part.AVD[padre].Prop
			j := 0
			for j = 0; j < len(part.BitmapAVD); j++ {
				if part.BitmapAVD[j] == 0 {
					if j == len(part.AVD) {
						part.AVD = append(part.AVD, avd)
					} else {
						part.AVD[j] = avd
					}
					part.BitmapAVD[j] = 1
					break
				}
			}
			if j == len(part.BitmapAVD) {
				fmt.Println("Error: ya ha creado la cantidad maxima de directorios")
				return 0, -1
			}
			part.AVD[padre].IndiceNext = int64(j) //le asigno el indice del anexo al next del padre
			i = 0
			padre = j
		}
	}
	return i, padre
}

func EncontrarCarpeta(particion *disco.Montada, nombre string, indice int) (int, int) { //busca si las carpetas de la ruta ya existen
	var name [20]byte // devolviendo el indice de la carpeta y su padre
	copy(name[:], nombre)
	index := 0
	padre := indice
	for i := 0; i < 6; i++ {
		index = int(particion.AVD[indice].IndicesSubs[i])
		if particion.AVD[index].Nombre == name {
			return index, padre
		}
	}
	index = int(particion.AVD[indice].IndiceNext)
	if index != 0 {
		index, padre = EncontrarCarpeta(particion, nombre, index)
	}
	return index, padre
}
