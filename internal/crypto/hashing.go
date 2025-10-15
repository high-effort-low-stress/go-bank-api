package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func HashTokenSHA256(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GenerateVerificationToken() (rawToken string, hashedToken string, err error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}

	rawToken = hex.EncodeToString(bytes)
	hashedToken = HashTokenSHA256(rawToken)

	return rawToken, hashedToken, nil
}
