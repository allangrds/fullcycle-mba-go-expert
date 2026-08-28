// Package httphandler traduz requisições HTTP em chamadas às regras de
// negócio de resolução de CEP e consulta de temperatura.
package httphandler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/6-clima-google-cloud-run/internal/cep"
	"github.com/allangrds/fullcycle-mba-go-expert/desafios/6-clima-google-cloud-run/internal/temperature"
	"github.com/go-chi/chi/v5"
)

// CityFinder resolve a cidade correspondente a um CEP.
type CityFinder interface {
	FindCity(ctx context.Context, cep string) (string, error)
}

// TemperatureFetcher consulta a temperatura atual (em Celsius) de uma cidade.
type TemperatureFetcher interface {
	CurrentCelsius(ctx context.Context, city string) (float64, error)
}

// Handler resolve o clima atual a partir de um CEP.
type Handler struct {
	Cities CityFinder
	Temps  TemperatureFetcher
}

type errorResponse struct {
	Message string `json:"message"`
}

// ServeHTTP implementa o contrato do desafio:
//   - 422 "invalid zipcode" se o CEP não tiver 8 dígitos.
//   - 404 "can not find zipcode" se o CEP não existir.
//   - 200 com temp_C/temp_F/temp_K em caso de sucesso.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zipcode := chi.URLParam(r, "cep")

	if !cep.IsValid(zipcode) {
		respondError(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	city, err := h.Cities.FindCity(r.Context(), zipcode)
	if err != nil {
		if errors.Is(err, cep.ErrNotFound) {
			respondError(w, http.StatusNotFound, "can not find zipcode")
			return
		}
		log.Printf("httphandler: failed to resolve city for cep %s: %v", zipcode, err)
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	celsius, err := h.Temps.CurrentCelsius(r.Context(), city)
	if err != nil {
		log.Printf("httphandler: failed to fetch temperature for city %s: %v", city, err)
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(temperature.FromCelsius(celsius))
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Message: message})
}
