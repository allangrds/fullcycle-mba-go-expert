package main

import "github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/configs"

func main() {
	config, _ := configs.LoadConfig(".")

	println(config.DBDriver)
}
