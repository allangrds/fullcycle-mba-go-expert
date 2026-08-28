package httphandler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/7-clima-cep-observabilidade/internal/viacep"
	"github.com/stretchr/testify/assert"
)

type fakeCityFinder struct {
	city string
	err  error
}

func (f fakeCityFinder) FindCity(ctx context.Context, cep string) (string, error) {
	return f.city, f.err
}

type fakeTemperatureFetcher struct {
	tempC float64
	err   error
}

func (f fakeTemperatureFetcher) CurrentCelsius(ctx context.Context, city string) (float64, error) {
	return f.tempC, f.err
}

func TestWeatherHandler_ServeHTTP(t *testing.T) {
	t.Run("invalid cep", func(t *testing.T) {
		handler := WeatherHandler{}
		req := httptest.NewRequest(http.MethodPost, "/weather", bytes.NewReader([]byte(`{"cep":"123"}`)))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.JSONEq(t, `{"message":"invalid zipcode"}`, rec.Body.String())
	})

	t.Run("cep not found", func(t *testing.T) {
		handler := WeatherHandler{
			Cities: fakeCityFinder{err: viacep.ErrNotFound},
			Temps:  fakeTemperatureFetcher{},
		}
		req := httptest.NewRequest(http.MethodPost, "/weather", bytes.NewReader([]byte(`{"cep":"99999999"}`)))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.JSONEq(t, `{"message":"can not find zipcode"}`, rec.Body.String())
	})

	t.Run("success", func(t *testing.T) {
		handler := WeatherHandler{
			Cities: fakeCityFinder{city: "São Paulo"},
			Temps:  fakeTemperatureFetcher{tempC: 28.5},
		}
		req := httptest.NewRequest(http.MethodPost, "/weather", bytes.NewReader([]byte(`{"cep":"29902555"}`)))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"city":"São Paulo","temp_C":28.5,"temp_F":83.3,"temp_K":301.5}`, rec.Body.String())
	})

	t.Run("temperature fetch error", func(t *testing.T) {
		handler := WeatherHandler{
			Cities: fakeCityFinder{city: "São Paulo"},
			Temps:  fakeTemperatureFetcher{err: errors.New("boom")},
		}
		req := httptest.NewRequest(http.MethodPost, "/weather", bytes.NewReader([]byte(`{"cep":"29902555"}`)))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
