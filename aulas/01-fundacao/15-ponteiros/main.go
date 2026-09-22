package main

func main() {
  a := 10

  //asterisco de ponteiro
  //o ponteiro aponta para o endereço de memória da variável a
  //& é o operador que retorna o endereço de memória da variável
  var ponteiro *int = &a

  println("Valor de a:", a)
  println("Endereço de a:", &a)
  println("Valor do ponteiro:", ponteiro)

  *ponteiro = 20 // altera o valor da variável através do ponteiro
  println("------")
  println("Valor de a:", a)
  println("Endereço de a:", &a)
  println("Valor do ponteiro:", ponteiro)

  b := &a
  println("------")
  println("Valor de b:", b)
  println("Valor real de b:", *b) // dereferenciação do ponteiro

  *b = 30
  println("------")
  println("Valor de a:", a)
  println("Endereço de a:", &a)
  println("Valor de b:", b)
  println("Valor real de b:", *b)
  println("Valor do ponteiro:", ponteiro)
  println("Valor real do ponteiro:", *ponteiro)
}
