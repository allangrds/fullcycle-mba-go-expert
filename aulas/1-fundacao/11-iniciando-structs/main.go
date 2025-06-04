package main

import "fmt"

type Client struct {
  Nome  string
  Idade int
  Ativo bool
}

func main() {
  cliente := Client{
    Nome: "João",
    Idade: 30,
    Ativo: true,
  }
  cliente.Ativo = false

  fmt.Printf("Nome: %s, Idade: %d, Ativo: %t\n", cliente.Nome, cliente.Idade, cliente.Ativo)
}
