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

var cantMontajes int64

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
					boot := crearSuperB(particion.Size, particion.Start, particion.Nombre)

					//EscribirMKFS(path, particion.Start)
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

func EscribirMKFS(path string, parte int64, superboot estructuras.SBoot) {
	f, err := os.OpenFile(path, os.O_RDWR, 777)
	if err != nil {
		log.Fatal(err)
	}
	f.Seek(parte, 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, superboot)
	disco.EscribirBytes(f, binario.Bytes())
	f.Close()
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
	t := time.Now()
	s := string(t.Format("Mon Jan _2 15:04:05 2006"))
	copy(superboot.Creacion[:], s)
	copy(superboot.LastMontaje[:], s)
	superboot.Montajes = cantMontajes
	cantMontajes++
	superboot.DetalleDirectorio = no
	superboot.ArbolVirtual = no
	superboot.Inodos = 5 * no
	superboot.SizeArbol = av
	superboot.SizeBloque = b
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
	superboot.Bitacora = superboot.InicioBloque + 20*no*b
	superboot.MagicNum = 201500332
	superboot.LibreArbolVirtual = no
	superboot.LibreBloques = 20 * no
	superboot.LibreDetalleDirec = no
	superboot.LibreInodos = 5 * no
	superboot.LibreBitInodo = superboot.BitmapTablaInodo
	superboot.LibreBitDetalle = superboot.BitmapDetalleDirec
	superboot.LibreBitBloque = superboot.BitmapBloques
	superboot.LibreBitArbol = superboot.BitmapArbol
	return superboot
}
