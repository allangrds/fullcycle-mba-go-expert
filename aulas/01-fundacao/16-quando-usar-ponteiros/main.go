package main

func soma(a, b int) int {
  return a + b
}

func soma2(a, b int) int {
  a = 50
  return a + b
}

func soma3(a *int, b int) int {
  *a = 50
  return *a + b
}

func main() {
  minhaVariavel := 1
  var minhaVariavel2 = 2

  println("Soma:", soma(minhaVariavel, minhaVariavel2))
  println("Soma 2:", soma2(minhaVariavel, minhaVariavel2))
  println("Minha Variável:", minhaVariavel)
  println("Minha Variável 2:", minhaVariavel2)
  println("Soma 3:", soma3(&minhaVariavel, minhaVariavel2))
  println("Minha Variável:", minhaVariavel)
  println("Minha Variável 2:", minhaVariavel2)
}
