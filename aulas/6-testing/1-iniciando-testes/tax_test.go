package tax

import "testing"

// Test<NomeDoTest> para criar o teste em Go
// <NomeDoText_test.go para criar o arquivo do teste em Go
// go test para rodar os testes
func TestCalculateTax(t *testing.T) {
	amount := 500.0
	expected := 5.0

	result := CalculateTax(amount)

	if result != expected {
		t.Errorf("Expected %f but got %f", expected, result)
	}
}
