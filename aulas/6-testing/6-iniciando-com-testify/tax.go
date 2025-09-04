package tax

import "errors"

func CalculateTax(amount float64) (float64, error) {
	if amount <= 0 {
		return 0, errors.New("amount must be greater than zero")
	}

	if amount >= 1000 && amount < 20000 {
		return 10.0, nil
	}

	if amount >= 20000 {
		return 20.0, nil
	}

	return 5.0, nil
}
