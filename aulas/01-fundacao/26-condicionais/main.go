package main

import "fmt"

func main() {
  var a int = 10
  var b int = 20

  if a > b {
    fmt.Println("A é maior que B")
  } else { //nao existe o else if ou elif
    fmt.Println("B é maior que A")
  }

  switch {
    case a > b:
      fmt.Println("A é maior que B")
    case a < b:
      fmt.Println("B é maior que A")
    default:
      fmt.Println("A é igual a B")
  }
}
