package main

import (
	"fmt"

	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/configs"
)

func main() {
	config := configs.Conf{
		DBDriver:      configs.GetDBDriver(),
		DBHost:        configs.GetDBHost(),
		DBPort:        configs.GetDBPort(),
		DBUser:        configs.GetDBUser(),
		DBPassword:    configs.GetDBPassword(),
		DBName:        configs.GetDBName(),
		WebServerPort: configs.GetWebServerPort(),
		JWTSecret:     configs.GetJWTSecret(),
		JWTExpiresIn:  configs.GetJWTExpiresIn(),
		TokenAuth:     configs.GetTokenAuth(),
	}

	fmt.Printf("DB Driver: %s\n", config.DBDriver)
	fmt.Printf("DB Host: %s\n", config.DBHost)
	fmt.Printf("DB Port: %s\n", config.DBPort)
	fmt.Printf("Web Server Port: %s\n", config.WebServerPort)
}
