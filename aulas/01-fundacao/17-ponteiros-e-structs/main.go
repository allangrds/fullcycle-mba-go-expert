package main

import "fmt"

type Cliente struct {
  nome string
}

func (cliente Cliente) andou() {
  cliente.nome = "Maria"
  fmt.Printf("O cliente %s andou\n", cliente.nome)
}

func (cliente *Cliente) andou2() {
  cliente.nome = "Maria"
  fmt.Printf("O cliente %s andou\n", cliente.nome)
}

// NewCliente cria uma nova instância de Cliente
// e retorna um ponteiro para ela.
// Isso é útil para evitar a criação de cópias desnecessárias
// e para permitir a modificação do objeto original.
func NewCliente() *Cliente {
  // Retorna um ponteiro para uma nova instância de Cliente
  return &Cliente{nome: "João"}
}

func main() {
  cliente := Cliente{
    nome: "João",
  }
  cliente.andou()
  fmt.Printf("O valor da struct está com o nome %s\n", cliente.nome)
  cliente.andou2()
  fmt.Printf("O valor da struct está com o nome %s\n", cliente.nome)
  println("------")

  cliente2 := NewCliente()
  fmt.Printf("O valor da struct cliente2 está com o nome %s\n", cliente2.nome)
  cliente2.andou()
  fmt.Printf("O valor da struct cliente2 está com o nome %s\n", cliente2.nome)
}
