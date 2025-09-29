package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

var number uint64 = 0

func main() {
	// Usando Mutex para proteger a variável 'number'
	// m := sync.Mutex{}
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	m.Lock()
	// 	number++
	// 	m.Unlock()
	// 	w.Write([]byte(fmt.Sprintf("Número: %d", number)))
	// })

	// Usando Operações Atômicas para proteger a variável 'number'
	// Operações atômicas são mais performáticas que Mutex
	// Operações atômicas são limitadas a tipos primitivos (int, uint, etc)
	// Operações atômicas não servem para estruturas complexas
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&number, 1)
		w.Write([]byte(fmt.Sprintf("Número: %d", number)))
	})

	http.ListenAndServe(":8080", nil)
}
