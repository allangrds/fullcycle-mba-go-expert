package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	//metricas, telemetria
	ctx = context.WithValue(ctx, "token", "senha123")
	bookHotel(ctx)
}

func bookHotel(ctx context.Context) {
	token := ctx.Value("token")

	fmt.Println(token)
}
