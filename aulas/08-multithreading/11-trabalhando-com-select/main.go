package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Message struct {
	id  int64
	Msg string
}

func main() {
	c1 := make(chan Message)
	c2 := make(chan Message)
	var i int64 = 0

	//RabbitMQ
	go func() {
		for i < 5 {
			atomic.AddInt64(&i, 1)
			msg := Message{id: i, Msg: "Hello from RabbitMQ"}
			c1 <- msg
		}
	}()

	//Kafka
	go func() {
		atomic.AddInt64(&i, 1)
		msg := Message{id: i, Msg: "Hello from Kafka"}
		c2 <- msg
	}()

	for i := 0; i < 5; i++ {
		select {
		case msg1 := <-c1: //rabbitmq
			fmt.Printf("Received from RabbitMQ: %+v\n", msg1)
		case msg2 := <-c2: // kafka
			fmt.Printf("Received from Kafka: %+v\n", msg2)
		case <-time.After(time.Second * 3):
			fmt.Println("Timeout after 3 seconds")
			// default:
			// 	println("No channels ready")
		}
	}
}
