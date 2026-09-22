package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

func main() {
	httpClient := http.Client{}
	jsonVar := bytes.NewBuffer([]byte(`{"name": "allan"}`))
	resp, err := httpClient.Post("http://google.com", "application/json", jsonVar)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	io.CopyBuffer(os.Stdout, resp.Body, nil)
}
