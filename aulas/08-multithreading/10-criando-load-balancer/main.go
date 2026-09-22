package main

import (
	"fmt"
	"time"
)

func worker(workerId int, data chan int) {
	for n := range data {
		fmt.Printf("Worker %d recebeu %d\n", workerId, n)
		time.Sleep(time.Second)
	}
}

func main() {
	data := make(chan int)
	qtdWorkers := 2

	// conforme data fornece dados, workers processam
	for i := 1; i <= qtdWorkers; i++ {
		go worker(i, data)
	}

	// Simulando envio de dados
	for i := 0; i < 10; i++ {
		data <- i
	}
}
