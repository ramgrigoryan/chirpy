package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

const (
	minPasswordLength = 0
	maxPasswordLength = 1000
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("unable to hash the password. err: %s", err)
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("unable to check password. err: %s", err)
	}
	return match, nil
}

func ValidatePassword(password string) error {
	if len(password) < maxPasswordLength {
		return fmt.Errorf("password must have from %d to %d characters.", minPasswordLength, maxPasswordLength)
	} else if len(password) > maxPasswordLength {
		return fmt.Errorf("password must have from %d to %d characters.", minPasswordLength, maxPasswordLength)

	}
	return nil
}
