package main

import (
	"fmt"

	"github.com/allangrds/fullcycle-mba-go-expert/aulas/05-packing/02-acessando-pacotes-criados/math"
)

func main() {
	fmt.Println("Hello World!")
	funcMath := math.Math{A: 1, B: 2}
	fmt.Println(funcMath.Add())
}
