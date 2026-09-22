package main

import "fmt"

func main() {
  var minhaVar interface{} = "Olá, Mundo!"

  // A asserção de tipo é usada para converter uma interface vazia em um tipo específico
  println(minhaVar.(string))

  result, ok := minhaVar.(int)
  fmt.Printf("O resultado da asserção é: %v, sucesso: %v\n", result, ok)

  res2 := minhaVar.(int)
  fmt.Printf("O resultado da asserção é: %v\n", res2)
}
