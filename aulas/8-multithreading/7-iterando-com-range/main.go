package main

import "fmt"

// thread 1
func main() {
	ch := make(chan int)

	//vai rodar em background
	go publish(ch)

	//não vai rodar em background
	reader(ch)
}

// thread 2
func publish(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}

	//fecha o canal para evitar deadlock
	//quando o reader tentar ler mais dados
	// do que o publisher enviou
	close(ch)
}

// thread 3
func reader(ch chan int) {
	for n := range ch {
		fmt.Println("Recebido:", n)
	}
}
