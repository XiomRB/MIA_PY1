package disco

import (
	"Archivos/PY1/estructuras"
	"fmt"
	"log"
	"strconv"
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
	copy(name[:], comando.Name)
	mbr := LeerDisco(comando.Path)
	if mbr.Size > 0 {
		part := extraerPart(comando.Path, name, mbr.Particiones) //obtine la info de la particion que se va a montar
		if part.Estado {                                         //si encontro la particion
			indice := -1
			for i := 0; i < len(DiscosMontados); i++ { //busca el disco montado
				indice = i
				if DiscosMontados[i].Path == comando.Path && DiscosMontados[i].Estado { //si es igual la particion esta montada
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
					DiscosMontados[indice] = disco
					fmt.Println("Particion montada exitosamente")
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
						DiscosMontados[indice] = disco
					} else {
						disco.Particiones[i] = mont
						DiscosMontados[indice] = disco
					}
					fmt.Println("Particion montada exitosamente")
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
				fmt.Println("Particion montada exitosamente")
			}
		}
	}
}

func Desmontar(u Unmount) {
	for i := 0; i < len(u.Desmontadas); i++ {
		letra := EncontrarLetra(byte(u.Desmontadas[i][2]))
		num, err := strconv.Atoi(string(u.Desmontadas[i][3]))
		if err != nil {
			log.Fatal(err)
		}
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
	if len(DiscosMontados) == 0 {
		fmt.Println("Aun no existen particiones montadas")
		return
	}
	for i := 0; i < len(DiscosMontados); i++ {
		if verifDiscoMontado(i) {
			disco := DiscosMontados[i]
			for j := 0; j < len(disco.Particiones); j++ {
				name := MostrarInfoMontada(disco, j)
				if len(name) > 0 {
					fmt.Println("id -> vd"+string(asignarLetra(i))+strconv.Itoa(j+1), "  path -> ", disco.Path, "  name -> ", name)
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
			m.Name = name
			return m
		}
		if parts[i].Tipo == GetChar("e") {
			ext = parts[i]
		}
	}
	if ext.Tipo != GetChar("e") {
		return m
	}
	ebr := LeerEBR(path, ext.Start)
	for ebr.Next != -1 && ebr.Name != name {
		ebr = LeerEBR(path, ebr.Next)
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
		return string(disco.Particiones[indice].Nombre[:])
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
		return GetChar("a")
	case 1:
		return GetChar("b")
	case 2:
		return GetChar("c")
	case 3:
		return GetChar("d")
	case 4:
		return GetChar("e")
	case 5:
		return GetChar("f")
	case 6:
		return GetChar("g")
	case 7:
		return GetChar("h")
	case 8:
		return GetChar("i")
	case 9:
		return GetChar("j")
	case 10:
		return GetChar("k")
	case 11:
		return GetChar("l")
	case 12:
		return GetChar("m")
	case 13:
		return GetChar("n")
	case 14:
		return GetChar("o")
	case 15:
		return GetChar("p")
	case 16:
		return GetChar("q")
	case 17:
		return GetChar("r")
	case 18:
		return GetChar("d")
	case 19:
		return GetChar("t")
	case 20:
		return GetChar("u")
	case 21:
		return GetChar("v")
	case 22:
		return GetChar("w")
	case 23:
		return GetChar("x")
	case 24:
		return GetChar("y")
	case 25:
		return GetChar("z")
	default:
		return GetChar("0")
	}
}

func EncontrarLetra(l byte) int {
	switch l {
	case GetChar("a"):
		return 0
	case GetChar("b"):
		return 1
	case GetChar("c"):
		return 2
	case GetChar("d"):
		return 3
	case GetChar("e"):
		return 4
	case GetChar("f"):
		return 5
	case GetChar("g"):
		return 6
	case GetChar("h"):
		return 7
	case GetChar("i"):
		return 8
	case GetChar("j"):
		return 9
	case GetChar("k"):
		return 10
	case GetChar("l"):
		return 11
	case GetChar("m"):
		return 12
	case GetChar("n"):
		return 13
	case GetChar("o"):
		return 14
	case GetChar("p"):
		return 15
	case GetChar("q"):
		return 16
	case GetChar("r"):
		return 17
	case GetChar("s"):
		return 18
	case GetChar("t"):
		return 19
	case GetChar("u"):
		return 20
	case GetChar("v"):
		return 21
	case GetChar("w"):
		return 22
	case GetChar("x"):
		return 23
	case GetChar("y"):
		return 24
	case GetChar("z"):
		return 25
	default:
		return 100
	}
}
