package auth

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// CheckPassword membandingkan password dengan hash bcrypt.
func CheckPassword(hash, password string) error {
	// Kompatibilitas hash bcrypt dari PHP/CodeIgniter
	if strings.HasPrefix(hash, "$2y$") {
		hash = "$2a$" + hash[4:]
	}

	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}

// HashPassword membuat hash bcrypt baru.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}