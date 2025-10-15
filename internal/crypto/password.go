package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// Params define os parâmetros de configuração para o Argon2id.
type Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var owaspRecommendedParams = &Params{
	Memory:      64 * 1024, // 64 MB
	Iterations:  3,
	Parallelism: 4,
	SaltLength:  16,
	KeyLength:   32,
}

// CreateHash gera um hash Argon2id para uma senha.
// O resultado é uma string no formato PHC, pronta para ser salva no banco de dados.
func HashPassword(password string) (string, error) {
	// 1. Gerar um salt criptograficamente seguro.
	salt := make([]byte, owaspRecommendedParams.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 2. Gerar o hash usando Argon2id.
	hash := argon2.IDKey([]byte(password), salt, owaspRecommendedParams.Iterations, owaspRecommendedParams.Memory, owaspRecommendedParams.Parallelism, owaspRecommendedParams.KeyLength)

	// 3. Codificar o salt e o hash em Base64.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 4. Montar a string no formato PHC.
	// Formato: $argon2id$v=19$m=<memory>,t=<iterations>,p=<parallelism>$<salt>$<hash>
	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	fullHash := fmt.Sprintf(format, argon2.Version, owaspRecommendedParams.Memory, owaspRecommendedParams.Iterations, owaspRecommendedParams.Parallelism, b64Salt, b64Hash)

	return fullHash, nil
}
