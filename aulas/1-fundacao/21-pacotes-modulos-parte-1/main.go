package main

import (
	"fmt"
	"fullcycle-mba-go-expert/aulas/1-fundacao/21-pacotes-modulos-parte-1/matematica"
)

func main() {
  soma := matematica.Soma(matematica.A, 20)

  fmt.Println("O resultado da soma é:", soma)
}
