package main

import (
	"fmt"

	//rodar go mod edit -replace github.com/allangrds/fullcycle-mba-go-expert/aulas/05-packing/05-go-mod-replace/math=../math para funcionar sem o pacote ainda estar publicado
	//go mod tidy após a linha de cima
	"github.com/allangrds/fullcycle-mba-go-expert/aulas/05-packing/05-go-mod-replace/math"
)

func main() {
	m := math.Math{A: 1, B: 2}
	fmt.Println(m.Add())
}
