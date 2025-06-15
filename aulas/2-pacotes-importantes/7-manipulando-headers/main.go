package main

import "net/http"

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World 2!"))
}

func main() {
	http.HandleFunc("/", BuscaCepHandler)
	http.ListenAndServe(":8080", nil)
}
