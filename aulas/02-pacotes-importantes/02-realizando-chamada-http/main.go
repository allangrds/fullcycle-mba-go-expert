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
	result, err := io.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Status Code: %d\n", request.StatusCode)
	fmt.Printf("Response Body: %s\n", string(result))
	request.Body.Close()
}
