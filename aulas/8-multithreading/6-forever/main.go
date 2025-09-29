package main

// thread 1
func main() {
	forever := make(chan bool)

	// thread 2
	go func() {
		for i := 0; i < 10; i++ {
			println(i)
		}

		forever <- true
	}()

	//apenas isso vai dar deadlock, pois ninguem vai ler o canal
	// forever <- true

	<-forever
}
