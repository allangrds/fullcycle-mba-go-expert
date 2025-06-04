package main

import "fmt"

func main() {
  total := func() int {
    return soma(1, 2, 3, 4, 5)
  }()

  fmt.Println(total)
}

func soma(numeros ...int) int {
  total := 0

  for _, numero := range numeros {
    total += numero
  }

  return total
}
