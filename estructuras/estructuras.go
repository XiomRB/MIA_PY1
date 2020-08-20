package estructuras

//MBR para crear Disco
type MBR struct {
	Size        int64
	Creacion    [30]byte
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
