package validators_test

import (
	"testing"

	"github.com/high-effort-low-stress/go-bank-api/internal/utils/validators"
	"github.com/stretchr/testify/assert"
)

func TestValidatePasswordPattern(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"Valid password", "StrongPass123!", true},
		{"Valid with hyphen and underscore", "My-Password_123", true},
		{"Too short (7 chars)", "Ab1!Ab1", false},
		{"Minimum valid (8 chars)", "Ab1!Ab12", true},
		{"Maximum valid (64 chars)", "Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!", true},
		{"Too long (65 chars)", "Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!Aa1!1", false},
		{"Missing lowercase", "STRONGPASS123!", false},
		{"Missing uppercase", "strongpass123!", false},
		{"Missing number", "StrongPass!", false},
		{"Missing special char", "StrongPass123", false},
		{"Contains valid char (space)", "Strong Pass123!", true},
		{"Contains valid char (dot)", "Strong.Pass123!", true},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validators.ValidatePasswordPattern(tt.password)
			assert.Equal(t, tt.expected, got, "Password: %s", tt.password)
		})
	}
}
