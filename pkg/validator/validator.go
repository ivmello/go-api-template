package validator

import (
	"errors"
	"net/mail"
	"regexp"
)

var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrWeakPassword    = errors.New("password must be at least 6 characters with at least one number and one letter")
	ErrInvalidUUID     = errors.New("invalid UUID format")
	ErrContentRequired = errors.New("content is required")
	ErrContentTooLong  = errors.New("content exceeds maximum length")
)

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmail
	}
	return nil
}

// ValidatePassword checks if the password meets requirements
func ValidatePassword(password string) error {
	if len(password) < 6 {
		return ErrWeakPassword
	}

	// Check for at least one letter
	letterRegex := regexp.MustCompile(`[a-zA-Z]`)
	if !letterRegex.MatchString(password) {
		return ErrWeakPassword
	}

	// Check for at least one number
	numberRegex := regexp.MustCompile(`[0-9]`)
	if !numberRegex.MatchString(password) {
		return ErrWeakPassword
	}

	return nil
}

// ValidateUUID validates UUID format
func ValidateUUID(uuid string) error {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(uuid) {
		return ErrInvalidUUID
	}
	return nil
}

// ValidateContent validates message content
func ValidateContent(content string, maxLength int) error {
	if content == "" {
		return ErrContentRequired
	}
	
	if len(content) > maxLength {
		return ErrContentTooLong
	}
	
	return nil
}