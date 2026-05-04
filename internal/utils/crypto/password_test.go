package crypto_test

import (
	"strings"
	"testing"

	"github.com/high-effort-low-stress/go-bank-api/internal/utils/crypto"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	t.Run("Should generate a valid Argon2id hash string", func(t *testing.T) {
		password := "securePassword123"
		hash, err := crypto.HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hash)

		// O formato esperado é: $argon2id$v=19$m=65536,t=3,p=4$<salt>$<hash>
		parts := strings.Split(hash, "$")

		// O primeiro elemento é vazio pois a string começa com $
		assert.Equal(t, 6, len(parts), "Hash format should have 6 parts separated by $")
		assert.Equal(t, "argon2id", parts[1])
		assert.Equal(t, "v=19", parts[2])
		assert.Contains(t, parts[3], "m=65536")
		assert.Contains(t, parts[3], "t=3")
		assert.Contains(t, parts[3], "p=4")

		// Salt e Hash devem estar presentes (Base64)
		assert.NotEmpty(t, parts[4], "Salt part should not be empty")
		assert.NotEmpty(t, parts[5], "Key part should not be empty")
	})

	t.Run("Should produce different hashes for the same password due to random salt", func(t *testing.T) {
		password := "samePassword"

		hash1, err1 := crypto.HashPassword(password)
		hash2, err2 := crypto.HashPassword(password)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2, "Hashes of the same password should be unique due to salt")
	})

	t.Run("Should handle empty password", func(t *testing.T) {
		password := ""
		hash, err := crypto.HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.True(t, strings.HasPrefix(hash, "$argon2id$"))
	})
}
