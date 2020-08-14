package comandos

//MBR para crear Disco
type MBR struct {
	Size        int64
	Creacion    [30]rune
	Indice      int
	Particiones [4]Particion
}

//InfoPart indicara donde se encuentran las particiones
type InfoPart struct {
	Tipo   rune
	Name   [16]rune
	Inicio int
	Size   int
	Ajuste rune
}

//EBR contiene info de las particiones logicas
type EBR struct {
	Status rune
	Fit    rune
	Start  int
	Size   int
	Next   int
	Name   [16]rune
}

//Particion primaria o extendida
type Particion struct {
	Status rune
	Tipo   rune
	Fit    rune
	Start  int
	Size   int
	Name   [16]rune
}
