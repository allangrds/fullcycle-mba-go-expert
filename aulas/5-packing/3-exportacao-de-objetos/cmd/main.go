package main

import (
	"fmt"

	"github.com/allangrds/fullcycle-mba-go-expert/aulas/5-packing/3-exportacao-de-objetos/math"
)

func main() {
	// fmt.Println("Hello World!")
	// funcMath := math.Math{A: 1, B: 2}
	// fmt.Println(funcMath.Add())
	// fmt.Println(math.X)
	//undefined math.x pq ele está em minúsculo(não exportado)
	// fmt.Println(math.x)

	//Agora com um "construtor"
	funcMathB := math.NewMathB(1, 2)
	fmt.Println(funcMathB.AddB())
}
