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
	NameDisc            [16]byte
	NoArbolVirtual      int64
	NoDetalleDirectorio int64
	NoInodos            int64
	NoBloques           int64
	NoLibreArbolVirtual int64
	NoLibreDetalleDirec int64
	NoLibreInodos       int64
	NoLibreBloques      int64
	Creacion            [25]byte
	LastMontaje         [25]byte
	Montajes            int64
	BitmapArbol         int64
	InicioArbol         int64
	BitmapDetalleDirec  int64
	InicioDetalleDirec  int64
	BitmapTablaInodo    int64
	InicioInodo         int64
	BitmapBloques       int64
	InicioBloque        int64
	InicioBitacora      int64
	SizeArbol           int64
	SizeDetalleDirec    int64
	SizeInodo           int64
	SizeBloque          int64
	LibreBitArbol       int64
	LibreBitDetalle     int64
	LibreBitInodo       int64
	LibreBitBloque      int64
	MagicNum            int64
}

type AVD struct {
	Creacion       [25]byte
	Nombre         [20]byte
	Subdirectorios [6]int64
	DetalleDir     int64
	Next           int64
	Propietario    Usuario
	Permisos       int64
}

type DetalleDir struct {
	Files [5]File
	Next  int64
}

type File struct {
	Nombre   [20]byte
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
	Propietario Usuario
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

type Logueado struct {
	Name      [10]byte
	Estado    bool
	Grupo     [10]byte
	Particion string
}
type Usuario struct {
	Name   [10]byte
	Clave  [10]byte
	Estado bool
}

type Grupo struct {
	Indice   int
	Estado   bool
	Name     [10]byte
	Usuarios []Usuario
}
