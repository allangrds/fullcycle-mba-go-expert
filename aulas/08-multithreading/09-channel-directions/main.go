package main

import "fmt"

// func recebe(nome string, hello chan string) {
// hello chan<- indica que o canal é somente de receber
func recebe(nome string, hello chan<- string) {
	hello <- nome
}

// func ler(data chan string) {
// data <-chan indica que o canal é somente de enviar
func ler(data <-chan string) {
	fmt.Println(<-data)
}

// thread 1
func main() {
	hello := make(chan string)
	go recebe("hello", hello)
	ler(hello)
}
