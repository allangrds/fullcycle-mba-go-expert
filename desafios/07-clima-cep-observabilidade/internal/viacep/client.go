package viacep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
)

var ErrNotFound = errors.New("city not found for given cep")

var tracer = otel.Tracer("viacep-client")

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type viaCEPResponse struct {
	Localidade string `json:"localidade"`
}

func (c Client) FindCity(ctx context.Context, cep string) (string, error) {
	ctx, span := tracer.Start(ctx, "viacep-lookup")
	defer span.End()

	url := fmt.Sprintf("%s/ws/%s/json/", c.BaseURL, cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var body viaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}

	if body.Localidade == "" {
		return "", ErrNotFound
	}

	return body.Localidade, nil
}
