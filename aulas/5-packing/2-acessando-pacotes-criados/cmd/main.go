package main

import (
	"fmt"

	"github.com/allangrds/fullcycle-mba-go-expert/aulas/5-packing/2-acessando-pacotes-criados/math"
)

func main() {
	fmt.Println("Hello World!")
	funcMath := math.Math{A: 1, B: 2}
	fmt.Println(funcMath.Add())
}
