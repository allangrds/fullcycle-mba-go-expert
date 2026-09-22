package main

import (
	"fmt"
	"sync"
)

// thread 1
func main() {
	ch := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(10)

	//vai rodar em background
	go publish(ch)

	//vai rodar em background
	go reader(ch, &wg)

	wg.Wait()
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
func reader(ch chan int, wg *sync.WaitGroup) {
	for n := range ch {
		fmt.Println("Recebido:", n)
		wg.Done()
	}
}
