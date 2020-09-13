package reportes

import (
	"Archivos/PY1/comandos/disco"
	"Archivos/PY1/comandos/sistema"
	"fmt"
	"io/ioutil"
	"strings"
)

func selBitmap(comando Reporte) {
	letra, indice := sistema.EncontrarMontada(comando.Id)
	if letra == -1 {
		fmt.Println("Error: la particion no esta montada")
		return
	}
	switch strings.ToLower(comando.Name) {
	case "bm_arbdir":
		reportaBitmap(comando.Path, disco.DiscosMontados[letra].Particiones[indice].BitmapAVD)
	case "bm_detdir":
		reportaBitmap(comando.Path, disco.DiscosMontados[letra].Particiones[indice].BitmapDetalle)
	case "bm_inode":
		fmt.Println("Entre con inodo")
		reportaBitmap(comando.Path, disco.DiscosMontados[letra].Particiones[indice].BitmapInodo)
	case "bm_block":
		reportaBitmap(comando.Path, disco.DiscosMontados[letra].Particiones[indice].BitmapBloques)
	}
}

func reportaBitmap(path string, bitmap []byte) {
	var rep []byte
	f := obtenerCarpeta(path)
	salto := 1
	for i := 0; i < len(bitmap); i++ {
		if bitmap[i] == 0 {
			rep = append(rep, 48)
			rep = append(rep, 44)
		} else {
			rep = append(rep, 49)
			rep = append(rep, 44)
		}
		if salto%10 == 0 {
			rep = append(rep, 10)
		}
		salto++
	}
	ioutil.WriteFile(f, rep, 0755)
	fmt.Println("Reporte Creado")
}
