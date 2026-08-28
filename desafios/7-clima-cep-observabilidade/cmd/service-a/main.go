package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/7-clima-cep-observabilidade/internal/httphandler"
	otelprovider "github.com/allangrds/fullcycle-mba-go-expert/desafios/7-clima-cep-observabilidade/internal/platform/otel"
	"github.com/allangrds/fullcycle-mba-go-expert/desafios/7-clima-cep-observabilidade/internal/serviceaclient"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// serviceBHost é o nome do serviço B na rede do docker-compose: fixo pela
// topologia do compose, não varia por ambiente, então fica hardcoded aqui
// em vez de configurável — mesmo raciocínio usado para OTEL_SERVICE_NAME.
const serviceBHost = "service-b"

const requestTimeout = 5 * time.Second

func main() {
	ctx := context.Background()

	port := getEnv("SERVICE_A_PORT", "8080")
	serviceBURL := fmt.Sprintf("http://%s:%s", serviceBHost, getEnv("SERVICE_B_PORT", "8081"))
	collectorEndpoint := getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	serviceName := getEnv("OTEL_SERVICE_NAME", "service-a")

	shutdown, err := otelprovider.NewProvider(ctx, serviceName, collectorEndpoint)
	if err != nil {
		log.Fatalf("failed to init otel provider: %v", err)
	}
	defer shutdown(ctx)

	httpClient := &http.Client{
		Timeout:   requestTimeout,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	handler := httphandler.CEPForwardHandler{
		ServiceB: serviceaclient.Client{BaseURL: serviceBURL, HTTPClient: httpClient},
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Post("/cep", handler.ServeHTTP)
	router.Get("/health", healthCheckHandler)

	instrumentedRouter := otelhttp.NewHandler(router, serviceName, otelhttp.WithFilter(skipHealthCheck))

	log.Printf("service-a listening on :%s", port)
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
