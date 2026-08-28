// Package temperature converte uma temperatura em Celsius para Fahrenheit e Kelvin.
package temperature

import "math"

// Result é a temperatura nas três escalas exigidas pelo contrato da API.
type Result struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

// FromCelsius converte celsius para Fahrenheit e Kelvin, seguindo as
// fórmulas do CHALLANGE.md (F = C*1.8+32, K = C+273). Os três valores são
// arredondados para 2 casas decimais só para eliminar ruído de ponto
// flutuante (ex.: 28.5*1.8+32 = 83.30000000000001 em float64).
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
