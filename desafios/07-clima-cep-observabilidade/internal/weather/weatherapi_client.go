package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel"
)

var ErrWeatherAPI = fmt.Errorf("weatherapi request failed")

var tracer = otel.Tracer("weatherapi-client")

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

func (c Client) CurrentCelsius(ctx context.Context, city string) (float64, error) {
	ctx, span := tracer.Start(ctx, "weatherapi-lookup")
	defer span.End()

	endpoint := fmt.Sprintf("%s/v1/current.json?key=%s&q=%s", c.BaseURL, c.APIKey, url.QueryEscape(city))
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
		return 0, fmt.Errorf("%w: status %d", ErrWeatherAPI, resp.StatusCode)
	}

	var body weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, err
	}

	return body.Current.TempC, nil
}
