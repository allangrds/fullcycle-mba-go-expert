package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	request, err := http.Get("https://www.google.com")
	if err != nil {
		panic(err)
	}

	// O defer é usado para garantir que o método Close seja chamado
	// mesmo que ocorra um erro ou o código retorne antes.
	// Isso é útil para liberar recursos, como conexões de rede.
	// O defer é executado após a função main retornar, garantindo que o corpo da resposta seja fechado.
	defer request.Body.Close()

	result, err := io.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Status Code: %d\n", request.StatusCode)
	fmt.Printf("Response Body: %s\n", string(result))
}
