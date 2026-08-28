package viacep

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_FindCity(t *testing.T) {
	t.Run("returns city when found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"localidade":"São Paulo"}`))
		}))
		defer server.Close()

		client := Client{BaseURL: server.URL, HTTPClient: server.Client()}
		city, err := client.FindCity(context.Background(), "01001000")

		assert.NoError(t, err)
		assert.Equal(t, "São Paulo", city)
	})

	t.Run("returns ErrNotFound when localidade is empty", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"erro":"true"}`))
		}))
		defer server.Close()

		client := Client{BaseURL: server.URL, HTTPClient: server.Client()}
		_, err := client.FindCity(context.Background(), "99999999")

		assert.ErrorIs(t, err, ErrNotFound)
	})
}
