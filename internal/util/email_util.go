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

func IsOldEmailMatchNew(oldEmail string, newEmail string) (bool, error) {
	if oldEmail == newEmail {
		return true, errors.New("New Email is the same as old email")
	}
	return false, nil
}

func IsOldEmailMatchClaim(oldEmail string, claimEmail string) (bool, error) {
	if oldEmail != claimEmail {
		return false, errors.New("Old Email does not match claim email")
	}
	return true, nil
}
