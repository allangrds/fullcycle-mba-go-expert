package temperature

import "math"

type Result struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func FromCelsius(celsius float64) Result {
	return Result{
		Celsius:    round2(celsius),
		Fahrenheit: round2(celsius*1.8 + 32),
		Kelvin:     round2(celsius + 273),
	}
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
