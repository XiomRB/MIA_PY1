package disco

import (
	"Archivos/PY1/estructuras"
	"fmt"
)

type Montada struct { //particion
	Indice int
	Estado bool
	Nombre [16]byte
	Size   int64
	Start  int64
	Ajuste byte
}

type Montado struct { // disco
	Particiones []Montada
	Path        string
	Letra       byte
	Estado      bool
}

var DiscosMontados []Montado

type Mount struct {
	Path string
	Name string
}
type Unmount struct {
	Desmontadas []string
}

func Montar(comando Mount) {
	var name [16]byte
	for i := 0; i < 16; i++ {
		name[i] = comando.Name[i]
	}
	mbr := LeerDisco(comando.Path)
	if mbr.Size > 0 {
		part := extraerPart(comando.Path, name, mbr.Particiones) //obtine la info de la particion que se va a montar
		if part.Estado {                                         //si encontro la particion
			indice := -1
			for i := 0; i < len(DiscosMontados); i++ { //busca el disco montado
				indice = i
				if DiscosMontados[i].Path == comando.Path {
					break
				}
			}
			if indice == len(DiscosMontados) {
				for i := 0; i < indice; i++ {
					if !DiscosMontados[i].Estado {
						indice = i
						break
					}
				}
			}
			if indice > -1 && indice < len(DiscosMontados) { //si encontro el disco montado
				disco := DiscosMontados[indice]
				disco.Path = comando.Path
				disco.Estado = true
				var mont Montada
				if len(disco.Particiones) == 0 {
					mont = crearMontada(part, 1)
					disco.Particiones = append(disco.Particiones, mont)
				} else {
					for i := 0; i < len(disco.Particiones); i++ {
						mont = disco.Particiones[i]
						if mont.Nombre == name {
							fmt.Println("Error: La particion ya ha sido montada")
							return
						}
					}
					i := 0
					for i = 0; i < len(disco.Particiones); i++ {
						mont = disco.Particiones[i]
						if !mont.Estado {
							break
						}
					}
					mont = crearMontada(part, i+1)
					if i == len(disco.Particiones) {
						disco.Particiones = append(disco.Particiones, mont)
					} else {
						disco.Particiones[i] = mont
					}
				}
			} else { //sino montar disco y montar particion
				mont := crearMontada(part, 1)
				disco := Montado{}
				disco.Path = comando.Path
				disco.Estado = true
				if indice == -1 {
					indice = 0
				}
				disco.Letra = asignarLetra(indice)
				disco.Particiones = append(disco.Particiones, mont)
				DiscosMontados = append(DiscosMontados, disco)
			}
		}
	}
}

func Desmontar(u Unmount) {
	for i := 0; i < len(u.Desmontadas); i++ {
		letra := encontrarLetra(byte(u.Desmontadas[i][2]))
		num := int(u.Desmontadas[i][3])
		if len(DiscosMontados) > 0 {
			if verifDiscoMontado(letra) {
				disco := DiscosMontados[letra]
				if verifPartMontada(disco, num) {
					desm := Montada{}
					DiscosMontados[letra].Particiones[num-1] = desm
					if verificarMontadas(disco) {
						d := Montado{}
						DiscosMontados[letra] = d
					}
					fmt.Println("Particion desmontada con exito")
				} else {
					fmt.Println("Error: La particion no ha sido montada")
					return
				}
			} else {
				fmt.Println("Error: La particion no ha sido montada")
				return
			}
		} else {
			fmt.Println("Error: La particion no ha sido montada")
			return
		}
	}
}

func MostrarMontadas() {
	for i := 0; i < len(DiscosMontados); i++ {
		if verifDiscoMontado(i) {
			disco := DiscosMontados[i]
			for j := 0; j < len(disco.Particiones); j++ {
				name := MostrarInfoMontada(disco, j)
				if len(name) > 0 {
					fmt.Println("id -> vd", asignarLetra(i), j+1, "  path -> ", disco.Path, "  name -> ", name)
				}
			}

		}
	}
}

func extraerPart(path string, name [16]byte, parts [4]estructuras.Particion) estructuras.InfoPart {
	var ext estructuras.Particion
	var m estructuras.InfoPart
	m.Estado = false
	for i := 0; i < 4; i++ {
		if name == parts[i].Name {
			m.Size = parts[i].Size
			m.Inicio = parts[i].Start
			m.Estado = true
			m.Ajuste = parts[i].Fit
			return m
		}
		if parts[i].Tipo == getChar("e") {
			ext = parts[i]
		}
	}
	if ext.Tipo != getChar("e") {
		return m
	}
	ebr := leerEBR(path, ext.Start)
	for ebr.Next != -1 && ebr.Name != name {
		ebr = leerEBR(path, ebr.Next)
	}
	if ebr.Name == name {
		m.Estado = true
		m.Size = ebr.Size
		m.Inicio = ebr.Start
		m.Ajuste = ebr.Fit
	}
	return m
}

func MostrarInfoMontada(disco Montado, indice int) string {
	if verifPartMontada(disco, indice+1) {
		name := ""
		for i := 0; i < 16; i++ {
			name += string(disco.Particiones[indice].Nombre[i])
		}
		return name
	}
	return ""
}

func verificarMontadas(disco Montado) bool {
	for i := 0; i < len(disco.Particiones); i++ {
		if disco.Particiones[i].Estado {
			return false
		}
	}
	return true
}

func crearMontada(info estructuras.InfoPart, indice int) Montada {
	var particion Montada
	particion.Estado = true
	particion.Ajuste = info.Ajuste
	particion.Nombre = info.Name
	particion.Size = info.Size
	particion.Start = info.Inicio
	particion.Indice = indice
	return particion
}

func verifDiscoMontado(letra int) bool {
	if len(DiscosMontados) > letra {
		if DiscosMontados[letra].Estado {
			return true
		}
	}
	return false
}

func verifPartMontada(disco Montado, indice int) bool {
	if len(disco.Particiones) >= indice {
		if disco.Particiones[indice-1].Estado {
			return true
		}
	}
	return false
}

func asignarLetra(i int) byte {
	switch i {
	case 0:
		return getChar("a")
	case 1:
		return getChar("b")
	case 2:
		return getChar("c")
	case 3:
		return getChar("d")
	case 4:
		return getChar("e")
	case 5:
		return getChar("f")
	case 6:
		return getChar("g")
	case 7:
		return getChar("h")
	case 8:
		return getChar("i")
	case 9:
		return getChar("j")
	case 10:
		return getChar("k")
	case 11:
		return getChar("l")
	case 12:
		return getChar("m")
	case 13:
		return getChar("n")
	case 14:
		return getChar("o")
	case 15:
		return getChar("p")
	case 16:
		return getChar("q")
	case 17:
		return getChar("r")
	case 18:
		return getChar("d")
	case 19:
		return getChar("t")
	case 20:
		return getChar("u")
	case 21:
		return getChar("v")
	case 22:
		return getChar("w")
	case 23:
		return getChar("x")
	case 24:
		return getChar("y")
	case 25:
		return getChar("z")
	default:
		return getChar("0")
	}
}

func encontrarLetra(l byte) int {
	switch l {
	case getChar("a"):
		return 0
	case getChar("b"):
		return 1
	case getChar("c"):
		return 2
	case getChar("d"):
		return 3
	case getChar("e"):
		return 4
	case getChar("f"):
		return 5
	case getChar("g"):
		return 6
	case getChar("h"):
		return 7
	case getChar("i"):
		return 8
	case getChar("j"):
		return 9
	case getChar("k"):
		return 10
	case getChar("l"):
		return 11
	case getChar("m"):
		return 12
	case getChar("n"):
		return 13
	case getChar("o"):
		return 14
	case getChar("p"):
		return 15
	case getChar("q"):
		return 16
	case getChar("r"):
		return 17
	case getChar("s"):
		return 18
	case getChar("t"):
		return 19
	case getChar("u"):
		return 20
	case getChar("v"):
		return 21
	case getChar("w"):
		return 22
	case getChar("x"):
		return 23
	case getChar("y"):
		return 24
	case getChar("z"):
		return 25
	default:
		return 100
	}
}
