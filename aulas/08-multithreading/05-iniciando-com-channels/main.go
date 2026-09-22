package main

import "fmt"

// thread 1
func main() {
	//make cria um canal
	canal := make(chan string) //canal vazio

	// thread 2
	go func() {
		canal <- "Olá Mundo!" // canal cheio
	}()

	//thread 1
	mensagem := <-canal // canal vazio
	fmt.Println(mensagem)
}
