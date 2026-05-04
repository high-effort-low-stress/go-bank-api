package validators_test

import (
	"testing"

	"github.com/high-effort-low-stress/go-bank-api/internal/utils/validators"
	"github.com/stretchr/testify/assert"
)

func TestIsValidCPF(t *testing.T) {
	tests := []struct {
		name     string
		cpf      string
		expected bool
	}{
		{"Valid numeric CPF", "10088147096", true},
		{"Valid formatted CPF", "100.881.470-96", true},
		{"Invalid length", "123.456.789-0", false},
		{"All same digits", "111.111.111-11", false},
		{"Wrong first verifier digit", "064.510.000-40", false},
		{"Wrong second verifier digit", "064.510.000-31", false},
		{"Contains letters", "064.510.000-3A", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validators.IsValidCPF(tt.cpf)
			assert.Equal(t, tt.expected, got, "CPF: %s", tt.cpf)
		})
	}
}
