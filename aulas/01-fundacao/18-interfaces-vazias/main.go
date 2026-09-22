package main

import "fmt"

func showType(varType interface{}) {
  fmt.Printf("O tipo é: %T e o valor é %v\n", varType, varType)
}

func main() {
  var x interface{} = 10
  var y interface{} = "Hello, World!"

  showType(x)
  showType(y)
}
