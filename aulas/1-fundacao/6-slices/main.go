package main

import "fmt"

func main() {
  meuArray := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}
  fmt.Printf("len=%d cap=%d array=%v\n", len(meuArray), cap(meuArray), meuArray)
  fmt.Printf("len=%d cap=%d array=%v\n", len(meuArray[:0]), cap(meuArray[:0]), meuArray[:0])
  fmt.Printf("len=%d cap=%d array=%v\n", len(meuArray[:2]), cap(meuArray[:2]), meuArray[:2])
  fmt.Printf("len=%d cap=%d array=%v\n", len(meuArray[2:]), cap(meuArray[2:]), meuArray[2:])

  meuArray = append(meuArray, 22)
  fmt.Printf("len=%d cap=%d array=%v\n", len(meuArray[:2]), cap(meuArray[:2]), meuArray[:2])
  fmt.Printf("len=%d cap=%d array=%v\n", len(meuArray), cap(meuArray), meuArray)
}
