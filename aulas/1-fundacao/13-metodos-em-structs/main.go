package main

import "fmt"

type Endereco struct {
  Rua   string
  Numero int
  Cidade string
}

type Client struct {
  Nome     string
  Idade    int
  Ativo    bool
  // Endereco Endereco
  Endereco
}

func (cliente Client) Desativar() {
  cliente.Ativo = false
}

func main() {
  cliente := Client{
    Nome: "João",
    Idade: 30,
    Ativo: true,
  }
  cliente.Endereco = Endereco{
    Rua: "Rua das Flores",
    Numero: 123,
    Cidade: "São Paulo",
  }
  cliente.Endereco.Cidade = "Rio de Janeiro"
  cliente.Desativar()

  fmt.Printf("Nome: %s, Idade: %d, Ativo: %t\n", cliente.Nome, cliente.Idade, cliente.Ativo)
}
