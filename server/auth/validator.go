package auth

import (
	"fmt"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

const (
	minPasswordLength = 8
	maxPasswordLength = 128
	minEmailLength    = 3
	maxEmailLength    = 255
)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if len(email) < minEmailLength {
		return fmt.Errorf("email must be at least %d characters", minEmailLength)
	}

	if len(email) > maxEmailLength {
		return fmt.Errorf("email must not exceed %d characters", maxEmailLength)
	}

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) < minPasswordLength {
		return fmt.Errorf("password must be at least %d characters", minPasswordLength)
	}

	if len(password) > maxPasswordLength {
		return fmt.Errorf("password must not exceed %d characters", maxPasswordLength)
	}

	return nil
}
