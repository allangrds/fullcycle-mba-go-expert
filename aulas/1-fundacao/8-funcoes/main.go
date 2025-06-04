package main

import (
	"errors"
	"fmt"
)

// func soma(a int, b int) int {
//   return a + b
// }
func soma(a int, b int) (int, bool) {
  if a > 10 {
    return a + b, true
  }

  return a + b, false
}

func soma2(a, b int) (int, error) {
  if a > 10 {
    return 0, errors.New("a não pode ser maior que 10")
  }

  return a + b, nil
}

func main() {
  // println("A soma de", a, "e", b, "é:", soma(a, b))
  // fmt.Println(soma(1, 2))
  // fmt.Println(soma(11, 2))
  // fmt.Println(soma2(11, 2))
  // fmt.Println(soma2(9, 2))

  var resultadSoma2, resultadSoma2Err = soma2(9, 2)
  var resultadSoma21, resultadSoma21Err = soma2(11, 2)

  fmt.Println("Resultado da soma2:", resultadSoma2, "Erro:", resultadSoma2Err)
  fmt.Println("Resultado da soma21:", resultadSoma21, "Erro:", resultadSoma21Err)
}


