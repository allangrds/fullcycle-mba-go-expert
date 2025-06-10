package main

import "fmt"

func main() {
  //for comum
  for i:= 0; i < 10; i++ {
    fmt.Printf("Valor de i: %d\n", i)
  }

  //for range
  var numeros = []int{ 1, 2, 3, 4, 5 };
  for indice, valor := range numeros {
    fmt.Println("Indice:", indice, "| Valor:", valor)
  }

  //for range sem indice
  var numeros2 = []int{ 1, 2, 3, 4, 5 };
  for _, valor := range numeros2 {
    fmt.Println("Valor:", valor)
  }

  //for range sem valor
  var numeros3 = []int{ 1, 2, 3, 4, 5 };
  for key := range numeros3 {
    fmt.Println("Indice:", key)
  }

  //for com declaracao previa de contador
  i := 0
  for i < 10 {
    fmt.Printf("Valor de i: %d\n", i)
    i++
  }

  //for loop infinito
  for {
   fmt.Println("Loop infinito")
   break
  }
}
