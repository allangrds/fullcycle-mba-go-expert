package main

func main() {
	ch := make(chan string, 2)

	ch <- "Hello, World!"
	ch <- "Buffered Channel"

	println(<-ch)
}
