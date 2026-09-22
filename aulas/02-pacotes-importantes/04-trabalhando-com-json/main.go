package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Conta struct {
	Numero int
	Saldo  int
}

type Conta2 struct {
	Numero int `json:"n"`
	Saldo  int `json:"s"`
	// Usando o - para ignorar o campo Saldo na serialização JSON.
	// Isso significa que o campo Saldo não será incluído no JSON gerado.
	// Saldo  int `json:"-"`
}

func main() {
	conta := Conta{
		Numero: 12345,
		Saldo:  1000,
	}

	// Marshal converte a struct Conta em um JSON.
	// O método Marshal retorna um slice de bytes e um erro.
	res, err := json.Marshal(conta)
	if err != nil {
		panic(err)
	}

	println("Resultado em bytes:", res)
	println("Resultado em string:", string(res))
	fmt.Printf("Resultado em string: %s\n", res)

	//O NewEncoder cria um novo encoder que escreve para o os.Stdout.
	// Ele é usado para escrever dados JSON diretamente na saída padrão.
	// os.Stdout é o fluxo de saída padrão do terminal.
	// Isso é útil para imprimir dados JSON formatados no console.
	err = json.NewEncoder(os.Stdout).Encode(conta)
	if err != nil {
		panic(err)
	}

	var conta2 Conta
	// Unmarshal converte o JSON de volta para a struct Conta.
	// O método Unmarshal recebe um slice de bytes e preenche a struct.
	err = json.Unmarshal(res, &conta2)
	if err != nil {
		panic(err)
	}
	println(conta2.Saldo)

	jsonPuro := []byte(`{"n":2,"s":200}`)
	var conta3 Conta2
	err = json.Unmarshal(jsonPuro, &conta3)
	if err != nil {
		println(err)
	}
	println(conta3.Saldo)
}
