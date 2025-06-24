package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println("Request iniciado!")
	defer log.Println("Request finalizado!")

	select {
	case <-time.After(time.Second * 5):
		//Terminal, stdout
		log.Println("Request processada com sucesso!")
		//Retorno do servidor
		w.Write([]byte("Request processada com sucesso!"))
	case <-ctx.Done():
		log.Println("Request cancelada pelo client!")
		http.Error(w, "Request cancelada pelo client!", http.StatusRequestTimeout)
	}
}
