package cep

import "regexp"

var formatRegex = regexp.MustCompile(`^\d{8}$`)

func IsValid(cep string) bool {
	return formatRegex.MatchString(cep)
}
