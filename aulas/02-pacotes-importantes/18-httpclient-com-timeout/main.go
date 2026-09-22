package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	//Se mudar para time.Second a aplicação funciona
	httpClient := http.Client{Timeout: time.Microsecond}
	resp, err := httpClient.Get("http://google.com")
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
