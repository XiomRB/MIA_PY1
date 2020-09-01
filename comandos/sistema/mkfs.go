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
	"time"
	"unsafe"
)

type Mkfs struct {
	Id   string
	Tipo string
	Add  int64
	Unit string
}

func EncontrarMontada(id string) (disco.Montada, string) {
	letra := disco.EncontrarLetra(byte(id[2]))
	num, err := strconv.Atoi(string(id[3]))
	if err != nil {
		log.Fatal(err)
	}
	if disco.VerifDiscoMontado(letra) {
		if disco.VerifPartMontada(disco.DiscosMontados[letra], num) {
			return disco.DiscosMontados[letra].Particiones[num-1], disco.DiscosMontados[letra].Path
		}
	}
	vacia := disco.Montada{}
	return vacia, ""
}

func AdminComando(comando Mkfs) {
	if len(comando.Id) < 0 {
		fmt.Println("Error: El parametro id es obligatorio")
	} else {
		particion, path := EncontrarMontada(comando.Id)
		if len(comando.Tipo) > 0 {
			if comando.Add != 0 {
				fmt.Println("Error, el parametro id y el parametro tipo no pueden ser declarados juntos")
			} else {
				if len(path) != 0 {
					fmt.Println(particion)
					/*boot := crearSuperB(particion.Size, particion.Start, particion.Nombre)
					EscribirMKFS(path, particion.Start, boot)
					EscribirMKFS(path, particion.Start+int64(unsafe.Sizeof(boot)), boot)*/
				} else {
					fmt.Println("Error: la particion no ha sido montada")
				}
			}
		} else {
			if comando.Add == 0 {
				fmt.Println("Error: el parametro add recibe cualquier numero excepto 0")
			} else {
				if len(path) != 0 {

				} else {
					fmt.Println("Error: la particion no ha sido montada")
				}
			}
		}
	}
}

func creacionSistema(particion *disco.Montada) {
	particion.Superboot = crearSuperB(particion.Size, particion.Start, particion.Nombre)
	users := "1,G,root\n,1,U,root,root\n"
	usuario := crearUs("root", "201500332")
	particion.Grupos = append(particion.Grupos, estructuras.Grupo{})
	copy(particion.Grupos[0].Name[:], "root")
	particion.Grupos[0].Usuarios[0] = usuario
	bitmaparbol := make([]byte, particion.Superboot.NoArbolVirtual)
	bitmaparbol[0] = 1
	bitmapdetalle := make([]byte, particion.Superboot.NoDetalleDirectorio)
	bitmapdetalle[0] = 1
	bitmapinodo := make([]byte, particion.Superboot.NoInodos)
	bitmapinodo[0] = 1
	file := crearFile("users.txt", particion.Superboot.InicioInodo)
	detalle := estructuras.DetalleDir{}
	avd := crearAVD("/", usuario, 770)
	avd.DetalleDir = particion.Superboot.InicioDetalleDirec
	detalle.Files[0] = file
	inodo := CrearInodo(1, int64(len(users)))
	particion.Inodos = append(particion.Inodos, inodo)
	particion.DD = append(particion.DD, detalle)
	particion.AVD = append(particion.AVD)
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
	return us
}

/*func crearRoot(ind, bloque int64) {
	var escrito [33]byte
	copy(escrito[:], "1,G,root\n1,U,root,root,201500332\n")
	var inodo estructuras.Inodo
	inodo.Indice = 1
	inodo.Size = int64(len(escrito))
	inodo.NBloques = inodo.Size/25 + 1
	var bloques []estructuras.Bloque
	for i := int64(0); i < inodo.NBloques; i++ {
		bloques = append(bloques, estructuras.Bloque{})
		inodo.Bloques[i] = bloque + int64(25*i)
	}
	nb := 0
	bit := 0
	for i := 0; i < len(escrito); i++ {
		bloques[nb].Text[bit] = escrito[i]
		if i%24 == 0 {
			nb++
			bit = 0
		}
	}
	usuario := UsuarioLogueado("root", "root", "201500332")
	avd := crearAVD("/", usuario, 770)
	file := crearFile("users.txt", ind)
	detalle := estructuras.DetalleDir{}

}*/

func DarHora() string {
	t := time.Now()
	return string(t.Format("Mon Jan _2 15:04:05 2006"))
}

func crearAVD(name string, user estructuras.Usuario, permisos int64) estructuras.AVD {
	avd := estructuras.AVD{}
	copy(avd.Creacion[:], DarHora())
	copy(avd.Nombre[:], name)
	avd.Propietario = user
	avd.Permisos = permisos
	return avd
}

/*
func UsuarioLogueado(name, grup, clave string) estructuras.Usuario {
	usuario := estructuras.Usuario{}
	copy(usuario.Name[:], name)
	copy(usuario.Grupo[:], grup)
	copy(usuario.Clave[:], clave)
	usuario.Estado = true
	return usuario
}*/

func crearFile(name string, inodo int64) estructuras.File {
	file := estructuras.File{}
	copy(file.Nombre[:], "users.txt")
	copy(file.Creacion[:], DarHora())
	copy(file.Modif[:], DarHora())
	file.Inodo = inodo
	return file
}

func EscribirBloques(texto string, n int64) []estructuras.Bloque { ///asigna texto a los bloques
	bloques := make([]estructuras.Bloque, n)
	nb := 0
	c := 0
	for i := 0; i < len(texto); i++ {
		bloques[nb].Text[c] = byte(texto[i])
		c++
		if i%24 == 0 {
			c = 0
			nb++
		}
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
