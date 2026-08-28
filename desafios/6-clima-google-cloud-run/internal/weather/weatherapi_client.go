// Package weather consulta a temperatura atual de uma cidade.
package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Client consulta a API do WeatherAPI para obter a temperatura atual de uma cidade.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

type weatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// CurrentCelsius retorna a temperatura atual de city, em graus Celsius.
func (c Client) CurrentCelsius(ctx context.Context, city string) (float64, error) {
	query := url.Values{}
	query.Set("key", c.APIKey)
	query.Set("q", city)
	query.Set("aqi", "no")

	endpoint := fmt.Sprintf("%s/current.json?%s", c.BaseURL, query.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weatherapi: unexpected status %d", resp.StatusCode)
	}

	var data weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.Current.TempC, nil
}
