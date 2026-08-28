package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/6-clima-google-cloud-run/internal/cep"
	"github.com/allangrds/fullcycle-mba-go-expert/desafios/6-clima-google-cloud-run/internal/httphandler"
	"github.com/allangrds/fullcycle-mba-go-expert/desafios/6-clima-google-cloud-run/internal/weather"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	requestTimeout = 5 * time.Second
	viaCEPBaseURL  = "https://viacep.com.br"
	weatherAPIURL  = "https://api.weatherapi.com/v1"
)

func main() {
	port := getEnv("PORT", "8080")

	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		log.Fatal("WEATHER_API_KEY is required")
	}

	httpClient := &http.Client{Timeout: requestTimeout}

	handler := httphandler.Handler{
		Cities: cep.Client{
			BaseURL:    viaCEPBaseURL,
			HTTPClient: httpClient,
		},
		Temps: weather.Client{
			BaseURL:    weatherAPIURL,
			APIKey:     weatherAPIKey,
			HTTPClient: httpClient,
		},
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Get("/{cep}", handler.ServeHTTP)

	addr := ":" + port
	log.Printf("weather server running on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
