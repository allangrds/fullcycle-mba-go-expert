package main

import (
	"fmt"
	"sync"
	"time"
)

func task(name string, wg *sync.WaitGroup) {
	for i := 0; i < 5; i++ {
		fmt.Printf("%d Task %s is running\n", i, name)
		time.Sleep(1 * time.Second) // Simula um trabalho demorado
		wg.Done()
	}
}

// O main é a thread 1
// O garbage collector é a thread 2
func main() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(15)

	// Thread 3
	go task("A", &waitGroup)

	// Thread 4
	go task("B", &waitGroup)

	// Thread 5
	// Função anônima
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Printf("%d Task anonymous is running\n", i)
			time.Sleep(1 * time.Second) // Simula um trabalho demorado
			waitGroup.Done()
		}
	}()

	// Não mostrouu nada porque a main terminou antes das outras goroutines
	// <sem conteúdo>
	waitGroup.Wait()
}
