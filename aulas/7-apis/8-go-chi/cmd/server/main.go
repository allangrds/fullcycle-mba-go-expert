package main

import (
	// "github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/configs"

	"fmt"
	"net/http"

	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/configs"
	_ "github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/docs"
	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/internal/entity"
	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/internal/infra/database"
	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//@title Go Expert API Example
//@version 1.0
//@description This is a sample server for a Go Expert API Example.
//@termsOfService http://swagger.io/terms/

//@contact.name API Support
//@contact.url http://www.swagger.io/support
//@contact.email

//@license.name Apache 2.0
//@license.url http://www.apache.org/licenses/LICENSE-2.0.html

//@host localhost:8000
//@BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	configs := configs.Conf{
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

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDB := database.NewProduct(db)
	userDB := database.NewUser(db)
	productHandler := handlers.NewProductHandler(productDB)
	userHandler := handlers.NewUserHandler(userDB, configs.TokenAuth, configs.JWTExpiresIn)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/products", func(router chi.Router) {
		// Middleware para verificar o token JWT em todas as rotas deste grupo
		router.Use(jwtauth.Verifier(configs.TokenAuth))

		// Middleware que exige autenticação para todas as rotas dentro deste grupo
		router.Use(jwtauth.Authenticator)

		router.Post("/", productHandler.CreateProduct)
		router.Get("/", productHandler.GetProduct)
		router.Get("/", productHandler.GetProducts)
		router.Get("/{id}", productHandler.GetProduct)
		router.Put("/{id}", productHandler.UpdateProduct)
		router.Delete("/{id}", productHandler.DeleteProduct)
	})

	router.Post("/users", userHandler.CreateUser)
	router.Post("/users/generate-token", userHandler.GetJWT)

	router.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	http.ListenAndServe(":8000", router)
}

// Middleware to log incoming requests
// NOT USED, ONLY AS EXAMPLE
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request details
		fmt.Printf("Received request: %s %s\n", r.Method, r.URL.Path)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
