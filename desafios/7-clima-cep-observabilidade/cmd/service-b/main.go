package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/7-clima-cep-observabilidade/internal/httphandler"
	otelprovider "github.com/allangrds/fullcycle-mba-go-expert/desafios/7-clima-cep-observabilidade/internal/platform/otel"
	"github.com/allangrds/fullcycle-mba-go-expert/desafios/7-clima-cep-observabilidade/internal/viacep"
	"github.com/allangrds/fullcycle-mba-go-expert/desafios/7-clima-cep-observabilidade/internal/weather"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	requestTimeout = 5 * time.Second
	viaCEPBaseURL  = "https://viacep.com.br"
	weatherAPIURL  = "https://api.weatherapi.com"
)

func main() {
	ctx := context.Background()

	port := getEnv("SERVICE_B_PORT", "8081")
	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		log.Fatal("WEATHER_API_KEY is required")
	}
	collectorEndpoint := getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	serviceName := getEnv("OTEL_SERVICE_NAME", "service-b")

	shutdown, err := otelprovider.NewProvider(ctx, serviceName, collectorEndpoint)
	if err != nil {
		log.Fatalf("failed to init otel provider: %v", err)
	}
	defer shutdown(ctx)

	httpClient := &http.Client{
		Timeout:   requestTimeout,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	handler := httphandler.WeatherHandler{
		Cities: viacep.Client{BaseURL: viaCEPBaseURL, HTTPClient: httpClient},
		Temps:  weather.Client{BaseURL: weatherAPIURL, APIKey: weatherAPIKey, HTTPClient: httpClient},
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Post("/weather", handler.ServeHTTP)
	router.Get("/health", healthCheckHandler)

	instrumentedRouter := otelhttp.NewHandler(router, serviceName, otelhttp.WithFilter(skipHealthCheck))

	log.Printf("service-b listening on :%s", port)
	if err := http.ListenAndServe(":"+port, instrumentedRouter); err != nil {
		log.Fatal(err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func skipHealthCheck(r *http.Request) bool {
	return r.URL.Path != "/health"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
