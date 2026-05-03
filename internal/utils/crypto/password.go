package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const memory64MiB = 64 * 1024

type Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var owaspRecommendedParams = &Params{
	Memory:      memory64MiB,
	Iterations:  3,
	Parallelism: 4,
	SaltLength:  16,
	KeyLength:   32,
}

func HashPassword(password string) (string, error) {
	salt := make([]byte, owaspRecommendedParams.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, owaspRecommendedParams.Iterations, owaspRecommendedParams.Memory, owaspRecommendedParams.Parallelism, owaspRecommendedParams.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	fullHash := fmt.Sprintf(format, argon2.Version, owaspRecommendedParams.Memory, owaspRecommendedParams.Iterations, owaspRecommendedParams.Parallelism, b64Salt, b64Hash)

	return fullHash, nil
}
