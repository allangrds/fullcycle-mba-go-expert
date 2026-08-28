package weather_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/6-clima-google-cloud-run/internal/weather"
	"github.com/stretchr/testify/assert"
)

func TestClient_CurrentCelsius_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "dummy-key", r.URL.Query().Get("key"))
		assert.Equal(t, "São Paulo", r.URL.Query().Get("q"))

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"location":{"name":"Sao Paulo"},"current":{"temp_c":28.5}}`))
	}))
	defer server.Close()

	client := weather.Client{BaseURL: server.URL, APIKey: "dummy-key", HTTPClient: server.Client()}

	celsius, err := client.CurrentCelsius(context.Background(), "São Paulo")

	assert.NoError(t, err)
	assert.Equal(t, 28.5, celsius)
}

func TestClient_CurrentCelsius_UnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := weather.Client{BaseURL: server.URL, APIKey: "invalid", HTTPClient: server.Client()}

	_, err := client.CurrentCelsius(context.Background(), "São Paulo")

	assert.Error(t, err)
}
