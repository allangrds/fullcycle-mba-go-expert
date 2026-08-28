package cep_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/6-clima-google-cloud-run/internal/cep"
	"github.com/stretchr/testify/assert"
)

func TestClient_FindCity_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"cep":"01001-000","logradouro":"Praça da Sé","localidade":"São Paulo","uf":"SP"}`))
	}))
	defer server.Close()

	client := cep.Client{BaseURL: server.URL, HTTPClient: server.Client()}

	city, err := client.FindCity(context.Background(), "01001000")

	assert.NoError(t, err)
	assert.Equal(t, "São Paulo", city)
}

func TestClient_FindCity_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// O ViaCEP real devolve o campo "erro" como string, não bool.
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"erro":"true"}`))
	}))
	defer server.Close()

	client := cep.Client{BaseURL: server.URL, HTTPClient: server.Client()}

	_, err := client.FindCity(context.Background(), "00000000")

	assert.ErrorIs(t, err, cep.ErrNotFound)
}

func TestClient_FindCity_UnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := cep.Client{BaseURL: server.URL, HTTPClient: server.Client()}

	_, err := client.FindCity(context.Background(), "01001000")

	assert.Error(t, err)
	assert.NotErrorIs(t, err, cep.ErrNotFound)
}
