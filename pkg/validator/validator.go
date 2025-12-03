package validator

import (
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidatePassword(password string) bool {
	return len(password) >= 6
}

func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

