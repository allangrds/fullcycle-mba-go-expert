package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

func BuscaCep(cep string) (*ViaCEP, error) {
	resp, error := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if error != nil {
		return nil, error
	}

	defer resp.Body.Close()

	body, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}
	var c ViaCEP
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}

	return &c, nil
}

// ResponseWriter é uma interface que permite escrever a resposta HTTP
// Ela é usada para enviar dados de volta ao cliente, como status, cabeçalhos e corpo da resposta.
// Request é uma estrutura que contém informações sobre a solicitação HTTP recebida pelo servidor.
// Ela inclui detalhes como o método HTTP, URL, cabeçalhos e corpo da solicitação.
func BuscaCepHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cep, error := BuscaCep(cepParam)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// NewEncoder é uma função que cria um novo encoder JSON que escreve para o ResponseWriter.
	// Ele é usado para codificar estruturas Go em JSON e enviá-las como resposta HTTP.
	// Encode é um método que codifica a estrutura Go fornecida em JSON e escreve no encoder.
	// Neste caso, ele está escrevendo a estrutura ViaCEP no formato JSON para o ResponseWriter.
	// Poderia ser feito de forma alternativa usando json.Marshal e w.Write, mas o Encoder facilita o processo
	json.NewEncoder(w).Encode(cep)
}

func main() {
	http.HandleFunc("/", BuscaCepHandler)
	http.ListenAndServe(":8080", nil)
}
