package cep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ErrNotFound indica que o CEP tem formato válido, mas não existe na base do ViaCEP.
var ErrNotFound = errors.New("cep not found")

// Client consulta a API do ViaCEP para resolver a cidade de um CEP.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// A resposta de "não encontrado" do ViaCEP é {"erro":"true"} — o campo
// "erro" vem como string, não bool. Em vez de decodificar esse campo (e
// lidar com essa inconsistência de tipo), usamos a ausência de
// "localidade" como sinal de CEP inexistente: ela só vem vazia quando o
// CEP não existe.
type viaCEPResponse struct {
	Localidade string `json:"localidade"`
}

// FindCity retorna o nome da cidade associada ao cep, ou ErrNotFound se o
// CEP não existir na base do ViaCEP.
func (c Client) FindCity(ctx context.Context, cep string) (string, error) {
	url := fmt.Sprintf("%s/ws/%s/json/", c.BaseURL, cep)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("viacep: unexpected status %d", resp.StatusCode)
	}

	var data viaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if data.Localidade == "" {
		return "", ErrNotFound
	}

	return data.Localidade, nil
}
