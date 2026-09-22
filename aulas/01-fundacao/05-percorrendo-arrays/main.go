package main

import "fmt"

var (
  a string = "Olá, Mundo!"
  meuArray [3]int
  meuArray2 [3]int = [3]int{1, 2, 3}
)

func main() {
  meuArray[0] = 10
  meuArray[1] = 20
  meuArray[2] = 30

  fmt.Println(len(meuArray))
  fmt.Println(len(meuArray2))
  fmt.Println(len(meuArray2) - 1)
  fmt.Println(meuArray2[0])
  fmt.Println(meuArray2[2])
  fmt.Println(meuArray2[len(meuArray2)-1])

  for i, value := range meuArray2 {
    fmt.Printf("%v. O índice %d tem o valor %d \n", a, i, value)
  }
}
