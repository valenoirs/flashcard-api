package password

import (
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {
	cost int
}

func NewPasswordHasher(cost int) *PasswordHasher {
	return &PasswordHasher{
		cost: cost,
	}
}

// preHash addresses the 72-byte bcrypt limit by hashing the password
// with SHA-256 and base64 encoding the result before giving it to bcrypt.
func preHash(plain string) []byte {
	hash := sha256.Sum256([]byte(plain))
	b64 := base64.StdEncoding.EncodeToString(hash[:])
	return []byte(b64)
}

// Hash securely hashes a plain text password.
func (p *PasswordHasher) Hash(plain string, cost int) (string, error) {
	b, err := bcrypt.GenerateFromPassword(preHash(plain), cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Verify checks if the plain text password matches the hashed version.
// Returning an error allows the caller to distinguish between a wrong
// password and a malformed hash format.
func (p *PasswordHasher) Verify(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), preHash(plain))
}
