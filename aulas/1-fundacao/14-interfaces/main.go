package main

import "fmt"

type Endereco struct {
  Rua   string
  Numero int
  Cidade string
}

// Interface é aplicada automaticamente a structs que implementam os métodos definidos na interface
// Interface aceita apenas métodos públicos
// Não é necessário declarar que a struct implementa a interface
type Pessoa interface {
  Desativar()
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
  fmt.Printf("Cliente %s desativado \n", cliente.Nome)
}

func Desativacao(pessoa Pessoa) {
  pessoa.Desativar()
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
  // cliente.Desativar()
  Desativacao(cliente)

  fmt.Printf("Nome: %s, Idade: %d, Ativo: %t\n", cliente.Nome, cliente.Idade, cliente.Ativo)
}
