package util

import (
	"errors"
	"unicode"
)

// ValidateNewPassword uses bcrypt's byte limit and a minimum length for new passwords.
// Existing passwords remain valid for authentication until rotated.
func ValidateNewPassword(password string) error {
	if len(password) < 12 || len(password) > 72 {
		return errors.New("invalidPassword")
	}
	var letter, digit, symbol bool
	for _, r := range password {
		letter = letter || unicode.IsLetter(r)
		digit = digit || unicode.IsDigit(r)
		symbol = symbol || (!unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r))
	}
	if !letter || !digit || !symbol {
		return errors.New("invalidPassword")
	}
	return nil
}
