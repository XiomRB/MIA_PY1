package estructuras

//MBR para crear Disco
type MBR struct {
	Size        int64
	Creacion    [25]byte
	Indice      int64
	Particiones [4]Particion
}

//InfoPart indicara donde se encuentran las particiones
type InfoPart struct {
	Tipo   byte
	Name   [16]byte
	Inicio int64
	Size   int64
	Ajuste byte
	Estado bool
}

//EBR contiene info de las particiones logicas
type EBR struct {
	Status byte
	Fit    byte
	Start  int64
	Size   int64
	Next   int64
	Name   [16]byte
}

//Particion primaria o extendida
type Particion struct {
	Status byte
	Tipo   byte
	Fit    byte
	Start  int64
	Size   int64
	Name   [16]byte
}

type SBoot struct {
	NameDisc           [16]byte
	ArbolVirtual       int64
	DetalleDirectorio  int64
	Inodos             int64
	Bloques            int64
	LibreArbolVirtual  int64
	LibreDetalleDirec  int64
	LibreInodos        int64
	LibreBloques       int64
	Creacion           [25]byte
	LastMontaje        [25]byte
	Montajes           int64
	BitmapArbol        int64
	InicioArbol        int64
	BitmapDetalleDirec int64
	InicioDetalleDirec int64
	BitmapTablaInodo   int64
	InicioInodo        int64
	BitmapBloques      int64
	InicioBloque       int64
	Bitacora           int64
	SizeArbol          int64
	SizeDetalleDirec   int64
	SizeInodo          int64
	SizeBloque         int64
	LibreBitArbol      int64
	LibreBitDetalle    int64
	LibreBitInodo      int64
	LibreBitBloque     int64
	MagicNum           int64
}

type AVD struct {
	Creacion       [25]byte
	Nombre         [16]byte
	Subdirectorios [6]int64
	DetalleDir     int64
	Next           int64
	Propietario    [10]byte
}

type DetalleDir struct {
	Files [5]File
	Next  int64
}

type File struct {
	Nombre   [16]byte
	Inodo    int64
	Creacion [25]byte
	Modif    [25]byte
	//Permiso int8
}

type Inodo struct {
	Indice      int64
	Size        int64
	NBloques    int64
	Bloques     [4]int64
	Indirecto   int64
	Propietario [10]byte //ver si no se quita
}

type Bloque struct {
	Text [25]byte
}

type Log struct {
	Operacion byte
	Tipo      byte
	Nombre    [16]byte
	Contenido bool
	Fecha     [25]byte
}

type Usuario struct {
	name   [10]byte
	clave  [10]byte
	estado bool
}

type Grupo struct {
	indice   int
	estado   bool
	usuarios [10]Usuario
}
