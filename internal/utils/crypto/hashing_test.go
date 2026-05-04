package crypto_test

import (
	"testing"

	"github.com/high-effort-low-stress/go-bank-api/internal/utils/crypto"
	"github.com/stretchr/testify/assert"
)

func TestHashTokenSHA256(t *testing.T) {
	t.Run("Should return deterministic hash for a given token", func(t *testing.T) {
		token := "my-secret-token"
		// echo -n "my-secret-token" | sha256sum
		expectedHash := "ea5add57437cbf20af59034d7ed17968dcc56767b41965fcc5b376d45db8b4a3"

		result := crypto.HashTokenSHA256(token)

		assert.Equal(t, expectedHash, result)
		assert.Len(t, result, 64, "SHA256 hex string should be 64 characters long")
	})

	t.Run("Should return different hashes for different tokens", func(t *testing.T) {
		hash1 := crypto.HashTokenSHA256("token1")
		hash2 := crypto.HashTokenSHA256("token2")

		assert.NotEqual(t, hash1, hash2)
	})
}

func TestGenerateVerificationToken(t *testing.T) {
	t.Run("Should generate valid raw and hashed tokens", func(t *testing.T) {
		raw, hashed, err := crypto.GenerateVerificationToken()

		assert.NoError(t, err)
		assert.NotEmpty(t, raw)
		assert.NotEmpty(t, hashed)

		// Verify raw token is 32 bytes hex encoded (64 chars)
		assert.Len(t, raw, 64)
		// Verify the relationship: hashed should be SHA256(raw)
		assert.Equal(t, crypto.HashTokenSHA256(raw), hashed)
	})

	t.Run("Should generate unique tokens on subsequent calls", func(t *testing.T) {
		raw1, _, _ := crypto.GenerateVerificationToken()
		raw2, _, _ := crypto.GenerateVerificationToken()

		assert.NotEqual(t, raw1, raw2)
	})
}
