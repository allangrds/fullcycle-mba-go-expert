package main

import (
	"fmt"
	"time"
)

func task(name string) {
	for i := 0; i < 5; i++ {
		fmt.Printf("%d Task %s is running\n", i, name)
		time.Sleep(1 * time.Second) // Simula um trabalho demorado
	}
}

// O main é a thread 1
// O garbage collector é a thread 2
func main() {
	// Thread 3
	go task("A")

	// Thread 4
	go task("B")

	// Thread 5
	// Função anônima
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Printf("%d Task anonymous is running\n", i)
			time.Sleep(1 * time.Second) // Simula um trabalho demorado
		}
	}()

	// Não mostrouu nada porque a main terminou antes das outras goroutines
	// <sem conteúdo>

	// Com o Sleep, a main espera as outras goroutines terminarem
	// Espera 6 segundos para ver a execução das goroutines
	fmt.Println("Main function is done")
	time.Sleep(2 * time.Second)
}
