package util

import (
	"fmt"
	"regexp"
	"strings"
)

func ValidatePassword(password string) error {
	// Must be at least 8 characters long
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	// Must contain at least one uppercase letter
	if password == strings.ToLower(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	// Must contain at least one lowercase letter
	if password == strings.ToUpper(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	// Must contain at least one number
	re := regexp.MustCompile(`[0-9]`)
	if !re.MatchString(password) {
		return fmt.Errorf("password must contain at least one number")
	}
	// Must contain at least one special character
	re = regexp.MustCompile(`[^a-zA-Z0-9]`)
	if !re.MatchString(password) {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}
