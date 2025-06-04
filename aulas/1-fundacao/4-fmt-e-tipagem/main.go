package main

import "fmt"

type ID int

var (
  a int = 123
  b ID  = 456
)

func main() {
  fmt.Printf("O tipo de 'a' é %T \n", a)
  fmt.Printf("O valor de 'a' é %v \n", a)
  fmt.Printf("O tipo de 'b' é %T \n", b)
}
