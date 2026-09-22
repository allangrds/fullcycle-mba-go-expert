package cep

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		name string
		cep  string
		want bool
	}{
		{name: "valid 8 digits", cep: "29902555", want: true},
		{name: "less than 8 digits", cep: "2990255", want: false},
		{name: "more than 8 digits", cep: "299025555", want: false},
		{name: "contains letters", cep: "2990255a", want: false},
		{name: "empty string", cep: "", want: false},
		{name: "with dash", cep: "29902-555", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValid(tt.cep))
		})
	}
}
