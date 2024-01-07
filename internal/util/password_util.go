package util

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(passwd string) (string, error) {
	if passwd == "" {
		return "", errors.New("password cannot empty")
	}
	if len(passwd) < 4 {
		return "", errors.New("password cannot be less than 4 characters")
	}

	result, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func ValidateHash(secret, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
	return err == nil
}
