package main

import (
	"fmt"
	"net/http"
)

var number uint64 = 0

func main() {
	//Pra cada request é criada uma nova thread
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		number++
		w.Write([]byte(fmt.Sprintf("Número: %d", number)))
	})

	http.ListenAndServe(":8080", nil)
}
