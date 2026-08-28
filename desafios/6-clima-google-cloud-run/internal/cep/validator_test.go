package cep_test

import (
	"testing"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/6-clima-google-cloud-run/internal/cep"
	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	cases := []struct {
		name string
		cep  string
		want bool
	}{
		{"8 digitos", "01001000", true},
		{"vazio", "", false},
		{"menos de 8 digitos", "123456", false},
		{"mais de 8 digitos", "123456789", false},
		{"letras", "abcdefgh", false},
		{"com traco", "01001-000", false},
		{"com espaco", "0100100 ", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, cep.IsValid(tc.cep))
		})
	}
}
