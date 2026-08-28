package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_CurrentCelsius(t *testing.T) {
	t.Run("returns temperature on success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"current":{"temp_c":28.5}}`))
		}))
		defer server.Close()

		client := Client{BaseURL: server.URL, APIKey: "fake-key", HTTPClient: server.Client()}
		temp, err := client.CurrentCelsius(context.Background(), "São Paulo")

		assert.NoError(t, err)
		assert.Equal(t, 28.5, temp)
	})

	t.Run("returns error when status is not 200", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":{"message":"invalid api key"}}`))
		}))
		defer server.Close()

		client := Client{BaseURL: server.URL, APIKey: "bad-key", HTTPClient: server.Client()}
		_, err := client.CurrentCelsius(context.Background(), "São Paulo")

		assert.ErrorIs(t, err, ErrWeatherAPI)
	})
}
