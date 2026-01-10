package auth

import "fmt"

var (
	ErrUnauthorized = fmt.Errorf("unauthorized: authentication required")
)
