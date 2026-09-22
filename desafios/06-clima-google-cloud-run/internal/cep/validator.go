// Package cep resolve a cidade correspondente a um CEP brasileiro.
package cep

import "regexp"

var formatRegex = regexp.MustCompile(`^[0-9]{8}$`)

// IsValid reporta se cep tem exatamente 8 dígitos numéricos.
func IsValid(cep string) bool {
	return formatRegex.MatchString(cep)
}
