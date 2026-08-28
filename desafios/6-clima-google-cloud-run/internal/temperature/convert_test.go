package temperature_test

import (
	"testing"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/6-clima-google-cloud-run/internal/temperature"
	"github.com/stretchr/testify/assert"
)

func TestFromCelsius(t *testing.T) {
	cases := []struct {
		name    string
		celsius float64
		want    temperature.Result
	}{
		{
			name:    "exemplo do CHALLANGE.md",
			celsius: 28.5,
			want:    temperature.Result{Celsius: 28.5, Fahrenheit: 83.3, Kelvin: 301.5},
		},
		{
			name:    "temperatura negativa",
			celsius: -10,
			want:    temperature.Result{Celsius: -10, Fahrenheit: 14, Kelvin: 263},
		},
		{
			name:    "zero",
			celsius: 0,
			want:    temperature.Result{Celsius: 0, Fahrenheit: 32, Kelvin: 273},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, temperature.FromCelsius(tc.celsius))
		})
	}
}
