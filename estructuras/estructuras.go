package estructuras

//MBR para crear Disco
type MBR struct {
	Size        int
	Creacion    [30]byte
	Indice      int
	Particiones [4]Particion
}

//InfoPart indicara donde se encuentran las particiones
type InfoPart struct {
	Tipo   byte
	Name   [16]byte
	Inicio int
	Size   int
	Ajuste byte
}

//EBR contiene info de las particiones logicas
type EBR struct {
	Status byte
	Fit    byte
	Start  int
	Size   int
	Next   int
	Name   [16]byte
}

//Particion primaria o extendida
type Particion struct {
	Status byte
	Tipo   byte
	Fit    byte
	Start  int
	Size   int
	Name   [16]byte
}
