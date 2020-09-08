package sistema

import (
	"Archivos/PY1/comandos/disco"
	"Archivos/PY1/estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type Mkfs struct {
	Id   string
	Tipo string
	Add  int64
	Unit string
}

func EncontrarMontada(id string) (int, int) { //retorna puntero para que todo se vaya guardando en la particion montada dentro de la lista de discos
	letra := disco.EncontrarLetra(byte(id[2]))
	num, err := strconv.Atoi(string(id[3]))
	if err != nil {
		log.Fatal(err)
	}
	if disco.VerifDiscoMontado(letra) {
		if disco.VerifPartMontada(disco.DiscosMontados[letra], num) {
			return letra, num - 1
		}
	}
	return -1, -1
}

func AdminComando(comando Mkfs) {
	if len(comando.Id) < 0 {
		fmt.Println("Error: El parametro id es obligatorio")
	} else {
		letra, indice := EncontrarMontada(comando.Id)
		if len(comando.Tipo) > 0 {
			if comando.Add != 0 {
				fmt.Println("Error, el parametro id y el parametro tipo no pueden ser declarados juntos")
			} else {
				if letra != -1 {
					if strings.EqualFold(comando.Tipo, "full") {
						f, err := os.OpenFile(disco.DiscosMontados[letra].Path, os.O_RDWR, 0755)
						if err != nil {
							log.Fatal(err)
						}
						f.Seek(int64(disco.DiscosMontados[letra].Particiones[indice].Start), 0)
						var cero int8 = 0
						t := int64(unsafe.Sizeof(cero))
						var binario bytes.Buffer
						binary.Write(&binario, binary.BigEndian, &cero)
						for i := int64(0); i <= disco.DiscosMontados[letra].Particiones[indice].Size/t; i += t {
							disco.EscribirBytes(f, binario.Bytes())
						}
						f.Close()
					}
					creacionSistema(&disco.DiscosMontados[letra].Particiones[indice])
					fmt.Println(len(disco.DiscosMontados[letra].Particiones[indice].AVD))
					//disco.EscribirSB(path, disco.DiscosMontados[letra].Particiones[indice].Start, disco.DiscosMontados[letra].Particiones[indice].Superboot)
					//disco.EscribirSB(path, disco.DiscosMontados[letra].Particiones[indice].Start+int64(unsafe.Sizeof(disco.DiscosMontados[letra].Particiones[indice].Superboot)), disco.DiscosMontados[letra].Particiones[indice].Superboot)
					fmt.Println("Sistema de archivos creado")
				} else {
					fmt.Println("Error: la particion no ha sido montada")
				}
			}
		} else {
			if comando.Add == 0 {
				fmt.Println("Error: el parametro add recibe cualquier numero excepto 0")
			} else {
				if letra != -1 {
					fmt.Println("Error: No se pude agregar el espacio deseado a la particion")
				} else {
					fmt.Println("Error: la particion no ha sido montada")
				}
			}
		}
	}
}

func creacionSistema(particion *disco.Montada) {
	particion.Superboot = crearSuperB(particion.Size, particion.Start, particion.Nombre)
	users := "1,G,root\n,1,U,root,201500332\n"
	usuario := crearUs("root", "201500332")
	particion.Grupos = append(particion.Grupos, estructuras.Grupo{})
	copy(particion.Grupos[0].Name[:], "root")
	particion.Grupos[0].Estado = true
	particion.Grupos[0].Usuarios = append(particion.Grupos[0].Usuarios, usuario)
	bitmaparbol := make([]byte, particion.Superboot.NoArbolVirtual)
	bitmaparbol[0] = 1
	bitmapdetalle := make([]byte, particion.Superboot.NoDetalleDirectorio)
	bitmapdetalle[0] = 1
	bitmapinodo := make([]byte, particion.Superboot.NoInodos)
	bitmapinodo[0] = 1
	file := CrearFile("users.txt", 0)
	detalle := estructuras.DetalleDir{}
	detalle.Next = -1
	avd := CrearAVD("/")
	avd.Prop.Name = usuario.Name
	avd.Prop.Grupo = usuario.Name
	avd.IndiceDD = 0
	detalle.Files[0] = file
	inodo := CrearInodo(1, int64(len(users)))
	particion.Inodos = append(particion.Inodos, inodo)
	particion.DD = append(particion.DD, detalle)
	particion.AVD = append(particion.AVD, avd)
	bloques := EscribirBloques(users, inodo.NBloques)
	for i := 0; i < len(bloques); i++ {
		inodo.Bloques[i] = particion.Superboot.InicioBloque + int64(i*int(unsafe.Sizeof(bloques[i])))
		particion.BB = append(particion.BB, bloques[i])
	}
	bitmapbloque := make([]byte, particion.Superboot.NoBloques)
	bitmapbloque[0] = 1
	bitmapbloque[1] = 1
	particion.BitmapInodo = bitmapinodo
	particion.BitmapDetalle = bitmapdetalle
	particion.BitmapAVD = bitmaparbol
	particion.BitmapBloques = bitmapbloque
	//modificarSB
	particion.Superboot.NoLibreArbolVirtual--
	particion.Superboot.LibreBitArbol++
	particion.Superboot.LibreBitBloque = particion.Superboot.BitmapBloques + 2
	particion.Superboot.LibreBitDetalle++
	particion.Superboot.LibreBitInodo++
	particion.Superboot.NoLibreBloques = particion.Superboot.NoLibreBloques - 2
	particion.Superboot.NoLibreDetalleDirec--
	particion.Superboot.NoLibreInodos--
}

func CrearInodo(indice int64, size int64) estructuras.Inodo {
	inodo := estructuras.Inodo{}
	inodo.Indice = indice
	inodo.Size = size
	inodo.Indirecto = -1
	if inodo.Size%25 == 0 {
		inodo.NBloques = inodo.Size / 25
	} else {
		inodo.NBloques = inodo.Size/25 + 1
	}
	return inodo
}

func NumEstructuras(size int64) (int64, int64, int64, int64, int64, int64) {
	subu := estructuras.SBoot{}
	avd := estructuras.AVD{}
	dde := estructuras.DetalleDir{}
	in := estructuras.Inodo{}
	bb := estructuras.Bloque{}
	ll := estructuras.Log{}
	sb := int64(unsafe.Sizeof(subu))
	av := int64(unsafe.Sizeof(avd))
	dd := int64(unsafe.Sizeof(dde))
	i := int64(unsafe.Sizeof(in))
	b := int64(unsafe.Sizeof(bb))
	l := int64(unsafe.Sizeof(ll))
	return (size - 2*sb) / (27 + av + dd + (5 * i) + (20 * b) + l), sb, av, dd, i, b
}

func crearSuperB(size, inicio int64, name [16]byte) estructuras.SBoot {
	no, sb, av, dd, i, b := NumEstructuras(size)
	superboot := estructuras.SBoot{}
	superboot.NameDisc = name
	copy(superboot.Creacion[:], DarHora())
	copy(superboot.LastMontaje[:], DarHora())
	superboot.Montajes = superboot.Montajes + 1
	superboot.NoDetalleDirectorio = no
	superboot.NoArbolVirtual = no
	superboot.NoInodos = 5 * no
	superboot.SizeArbol = av
	superboot.NoBloques = b
	superboot.SizeDetalleDirec = dd
	superboot.SizeInodo = i
	superboot.SizeBloque = b
	superboot.BitmapArbol = inicio + 2*sb
	superboot.InicioArbol = inicio + 2*sb + no
	superboot.BitmapDetalleDirec = superboot.InicioArbol + no*av
	superboot.InicioDetalleDirec = no + superboot.BitmapDetalleDirec
	superboot.BitmapTablaInodo = superboot.InicioDetalleDirec + no*dd
	superboot.InicioInodo = superboot.BitmapTablaInodo + no*5
	superboot.BitmapBloques = superboot.InicioInodo + 5*no*i
	superboot.InicioBloque = superboot.BitmapBloques + 20*no
	superboot.InicioBitacora = superboot.InicioBloque + 20*no*b
	superboot.MagicNum = 201500332
	superboot.NoLibreArbolVirtual = no
	superboot.NoLibreBloques = 20 * no
	superboot.NoLibreDetalleDirec = no
	superboot.NoLibreInodos = 5 * no
	superboot.LibreBitInodo = superboot.BitmapTablaInodo
	superboot.LibreBitDetalle = superboot.BitmapDetalleDirec
	superboot.LibreBitBloque = superboot.BitmapBloques
	superboot.LibreBitArbol = superboot.BitmapArbol
	return superboot
}

func crearUs(name, pass string) estructuras.Usuario {
	us := estructuras.Usuario{}
	copy(us.Name[:], name)
	copy(us.Clave[:], pass)
	us.Estado = true
	return us
}

func DarHora() string {
	t := time.Now()
	return string(t.Format("Mon Jan _2 15:04:05 2006"))
}

func CrearAVD(name string) estructuras.AVD {
	avd := estructuras.AVD{}
	copy(avd.Creacion[:], DarHora())
	copy(avd.Nombre[:], name)
	avd.Prop.Name = LoginUs.Name
	avd.Prop.Grupo = LoginUs.Grupo
	avd.Permisos[0] = 6
	avd.Permisos[1] = 4
	avd.Permisos[2] = 4
	avd.IndiceDD = -1
	return avd
}

func CrearFile(name string, inodo int64) estructuras.File {
	file := estructuras.File{}
	copy(file.Nombre[:], name)
	copy(file.Creacion[:], DarHora())
	copy(file.Modif[:], DarHora())
	file.Inodo = inodo
	file.Permisos[0] = 6
	file.Permisos[1] = 4
	file.Permisos[2] = 4
	file.Prop = estructuras.Propietario{}
	file.Prop.Name = LoginUs.Name
	file.Prop.Grupo = LoginUs.Grupo
	return file
}

func EscribirBloques(texto string, n int64) []estructuras.Bloque { ///asigna texto a los bloques
	bloques := make([]estructuras.Bloque, n)
	nb := 0
	c := 0
	for i := 0; i < len(texto); i++ {
		if i%25 == 0 && i > 0 {
			c = 0
			nb++
		}
		bloques[nb].Text[c] = byte(texto[i])
		c++
	}
	return bloques
}

func EscribirFile(path string, file *estructuras.File, seek int64) {
	f, err := os.OpenFile(path, os.O_RDWR, 755)
	if err != nil {
		log.Fatal(err)
	}
	f.Seek(seek, 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, &file)
	disco.EscribirBytes(f, binario.Bytes())
	f.Close()
}

func EscribirBloque(path string, bloque *estructuras.Bloque, seek int64) { //los va a escribir a disco
	f, err := os.OpenFile(path, os.O_RDWR, 755)
	if err != nil {
		log.Fatal(err)
	}
	f.Seek(seek, 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, &bloque)
	disco.EscribirBytes(f, binario.Bytes())
	f.Close()
}

func EscribirCarpeta(path string, avd *estructuras.AVD, seek int64) {
	f, err := os.OpenFile(path, os.O_RDWR, 755)
	if err != nil {
		log.Fatal(err)
	}
	f.Seek(seek, 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, &avd)
	disco.EscribirBytes(f, binario.Bytes())
	f.Close()
}

func EscribirDetalle(path string, detalle *estructuras.DetalleDir, seek int64) {
	f, err := os.OpenFile(path, os.O_RDWR, 755)
	if err != nil {
		log.Fatal(err)
	}
	f.Seek(seek, 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, &detalle)
	disco.EscribirBytes(f, binario.Bytes())
	f.Close()
}

func EscribirInodo(path string, inodo *estructuras.Inodo, seek int64) {
	f, err := os.OpenFile(path, os.O_RDWR, 755)
	if err != nil {
		log.Fatal(err)
	}
	f.Seek(seek, 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, &inodo)
	disco.EscribirBytes(f, binario.Bytes())
	f.Close()
}
