package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	//Se mudar para *3 vai pro ctx.Done
	ctx, cancel := context.WithTimeout(ctx, time.Second*6)
	defer cancel()
	bookHotel(ctx)
}

func bookHotel(ctx context.Context) {
	select {
	case <-ctx.Done():
		fmt.Println("Hotel booking cancelado. Timeout")
		return
	case <-time.After(time.Second * 5):
		fmt.Println("Hotel booked")
	}
}
