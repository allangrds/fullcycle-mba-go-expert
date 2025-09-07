package main

import (
	// "github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/configs"

	"net/http"

	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/internal/entity"
	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/internal/infra/database"
	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/internal/infra/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// config := configs.Conf{
	// 	DBDriver:      configs.GetDBDriver(),
	// 	DBHost:        configs.GetDBHost(),
	// 	DBPort:        configs.GetDBPort(),
	// 	DBUser:        configs.GetDBUser(),
	// 	DBPassword:    configs.GetDBPassword(),
	// 	DBName:        configs.GetDBName(),
	// 	WebServerPort: configs.GetWebServerPort(),
	// 	JWTSecret:     configs.GetJWTSecret(),
	// 	JWTExpiresIn:  configs.GetJWTExpiresIn(),
	// 	TokenAuth:     configs.GetTokenAuth(),
	// }

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	http.HandleFunc("/products", productHandler.CreateProduct)
	http.ListenAndServe(":8000", nil)
}
