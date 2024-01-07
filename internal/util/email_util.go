package util

import (
	"errors"
	"net/mail"
)

func ValidateEmail(email string) (bool, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false, err
	}
	return true, nil
}

func IsOldEmailMatch(oldEmail string, newEmail string) (bool, error) {
	if oldEmail == newEmail {
		return true, errors.New("New Email is the same as old email")
	}
	return false, nil
}
