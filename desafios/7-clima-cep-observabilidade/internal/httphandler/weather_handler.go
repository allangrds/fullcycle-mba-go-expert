package httphandler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/7-clima-cep-observabilidade/internal/cep"
	"github.com/allangrds/fullcycle-mba-go-expert/desafios/7-clima-cep-observabilidade/internal/temperature"
)

type CityFinder interface {
	FindCity(ctx context.Context, cep string) (string, error)
}

type TemperatureFetcher interface {
	CurrentCelsius(ctx context.Context, city string) (float64, error)
}

type WeatherHandler struct {
	Cities CityFinder
	Temps  TemperatureFetcher
}

type cepRequest struct {
	CEP string `json:"cep"`
}

type weatherResponse struct {
	City string `json:"city"`
	temperature.Result
}

func (h WeatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req cepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !cep.IsValid(req.CEP) {
		respondError(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	city, err := h.Cities.FindCity(r.Context(), req.CEP)
	if err != nil {
		respondError(w, http.StatusNotFound, "can not find zipcode")
		return
	}

	tempC, err := h.Temps.CurrentCelsius(r.Context(), city)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weatherResponse{
		City:   city,
		Result: temperature.FromCelsius(tempC),
	})
}
