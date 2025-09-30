package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var SERVER_TIMEOUT = 1 * time.Second

func makeGetRequest(url string) (string, error) {
	// Criando um contexto que será cancelado após 1 segundo
	createdContext, cancel := context.WithTimeout(context.Background(), SERVER_TIMEOUT)
	defer cancel() // Importante: sempre chame cancel para liberar recursos

	// Criando uma nova solicitação GET
	req, err := http.NewRequestWithContext(createdContext, "GET", url, nil)
	if err != nil {
		fmt.Println("Erro ao criar a solicitação:", err)
		return "", err
	}

	// Enviando a request HTTP
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro ao enviar request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Lida a resposta do servidor
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erro ao ler resposta do servidor:", err)
		return "", err
	}

	cepInfo := string(body)

	return cepInfo, nil
}

type CEPInfo struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type BrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

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

type APIResult struct {
	CEPInfo CEPInfo
	Source  string
}

func main() {
	// Criando channels para cada API
	brasilAPIChannel := make(chan APIResult)
	viaCEPChannel := make(chan APIResult)

	// Goroutine para BrasilAPI
	go func() {
		//Simulando atraso na resposta da BrasilAPI
		// time.Sleep(2 * time.Second)

		var FIRST_SERVER_API_URL = "https://brasilapi.com.br/api/cep/v1/01153000"

		// Criação da request
		response, err := makeGetRequest(FIRST_SERVER_API_URL)
		if err != nil {
			log.Printf("Erro ao obter resposta da BrasilAPI: %v", err)
			return
		}

		var brasilApi BrasilAPI
		err = json.Unmarshal([]byte(response), &brasilApi)
		if err != nil {
			log.Printf("Erro ao fazer unmarshal do JSON da BrasilAPI: %v", err)
			return
		}

		// Convertendo BrasilAPI para CEPInfo
		finalCep := CEPInfo(brasilApi)
		brasilAPIChannel <- APIResult{CEPInfo: finalCep, Source: "BrasilAPI"}
	}()

	// Goroutine para ViaCEP
	go func() {
		// Simulando atraso na resposta da ViaCEP
		// time.Sleep(2 * time.Second)

		var SECOND_SERVER_API_URL = "http://viacep.com.br/ws/01153000/json/"

		// Criação da request
		response, err := makeGetRequest(SECOND_SERVER_API_URL)
		if err != nil {
			log.Printf("Erro ao obter resposta da BrasilAPI: %v", err)
			return
		}

		var viaCep ViaCEP
		err = json.Unmarshal([]byte(response), &viaCep)
		if err != nil {
			log.Printf("Erro ao fazer unmarshal do JSON da VIACep: %v", err)
			return
		}

		// Convertendo BrasilAPI para CEPInfo
		finalCep := CEPInfo{
			Cep:          viaCep.Cep,
			State:        viaCep.Uf,
			City:         viaCep.Localidade,
			Neighborhood: viaCep.Bairro,
			Street:       viaCep.Logradouro,
		}
		viaCEPChannel <- APIResult{CEPInfo: finalCep, Source: "ViaCEP"}
	}()

	// Select para pegar a primeira resposta ou timeout
	select {
	case result := <-brasilAPIChannel:
		fmt.Printf("Resposta mais rápida veio da: %s\n", result.Source)
		fmt.Printf("CEP: %s\n", result.CEPInfo.Cep)
		fmt.Printf("Estado: %s\n", result.CEPInfo.State)
		fmt.Printf("Cidade: %s\n", result.CEPInfo.City)
		fmt.Printf("Bairro: %s\n", result.CEPInfo.Neighborhood)
		fmt.Printf("Rua: %s\n", result.CEPInfo.Street)
	case result := <-viaCEPChannel:
		fmt.Printf("Resposta mais rápida veio da: %s\n", result.Source)
		fmt.Printf("CEP: %s\n", result.CEPInfo.Cep)
		fmt.Printf("Estado: %s\n", result.CEPInfo.State)
		fmt.Printf("Cidade: %s\n", result.CEPInfo.City)
		fmt.Printf("Bairro: %s\n", result.CEPInfo.Neighborhood)
		fmt.Printf("Rua: %s\n", result.CEPInfo.Street)
	case <-time.After(time.Second * 1):
		fmt.Println("Timeout: nenhuma API respondeu dentro de 1 segundo")
	}
}
