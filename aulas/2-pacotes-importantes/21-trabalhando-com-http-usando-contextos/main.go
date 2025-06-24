package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	//Criação de um novo contexto
	ctx := context.Background()

	//Se o contexto passar de microsegundo, será cancelado
	ctx, cancel := context.WithTimeout(ctx, time.Microsecond)

	//Será cancelado de qualquer forma ao final, mesmo sem o timeout
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://google.com", nil)
	if err != nil {
		panic(err)
	}

	//Não preciso criar o http.Client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	println(string(body))
}
