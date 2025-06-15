package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	// os é um pacote que fornece uma interface para interagir com o sistema operacional
	// os.Args é um slice de strings que contém os argumentos da linha de comando
	// o 1: é para ignorar o primeiro argumento, que é o nome do programa
	for _, cep := range os.Args[1:] {
		var url = "http://viacep.com.br/ws/" + cep + "/json/"

		req, err := http.Get(url)
		if err != nil {
			//Estou usando o Fprintf que recebe como primeiro argumento o Writer, que nesse caso é o Stderr
			//O Stderr é um pacote que fornece uma interface para escrever na saída de erro padrão
			fmt.Fprintf(os.Stderr, "Erro ao fazer a requisição: %v\n", err)
		}
		defer req.Body.Close()

		res, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %v\n", err)
		}

		var data ViaCEP
		err = json.Unmarshal(res, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer parse da response: %v\n", err)
		}

		file, err := os.Create("cep.txt")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao criar arquivo: %v\n", err)
		}
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf("CEP: %s, Logradouro: %s\n", data.Cep, data.Logradouro))

	}
}
