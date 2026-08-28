package temperature

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromCelsius(t *testing.T) {
	tests := []struct {
		name    string
		celsius float64
		want    Result
	}{
		{
			name:    "zero celsius",
			celsius: 0,
			want:    Result{Celsius: 0, Fahrenheit: 32, Kelvin: 273},
		},
		{
			name:    "positive celsius",
			celsius: 28.5,
			want:    Result{Celsius: 28.5, Fahrenheit: 83.3, Kelvin: 301.5},
		},
		{
			name:    "negative celsius",
			celsius: -10,
			want:    Result{Celsius: -10, Fahrenheit: 14, Kelvin: 263},
		},
		{
			name:    "rounds to two decimals",
			celsius: 21.333,
			want:    Result{Celsius: 21.33, Fahrenheit: 70.4, Kelvin: 294.33},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromCelsius(tt.celsius)
			assert.Equal(t, tt.want, got)
		})
	}
}
