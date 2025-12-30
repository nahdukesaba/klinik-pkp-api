package utils

import "regexp"

// STRUCT TO INDICATE VALIDATOR UTILITIES
type Validator struct{}

// VALIDATE EMAIL FORMAT
func (v *Validator) ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// VALIDATE PASSWORD STRENGTH (MINIMUM 8 CHARACTERS)
func (v *Validator) ValidatePasswordStrength(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters"
	}
	return true, ""
}

// ValidatePhone validates phone number format
func (v *Validator) ValidatePhone(phone string) bool {
	if phone == "" {
		return true // Phone is optional
	}
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	return phoneRegex.MatchString(phone)
}

// ValidateRequired checks if value is not empty
func (v *Validator) ValidateRequired(value any, fieldName string) (bool, string) {
	if value == "" {
		return false, fieldName + " is required"
	}
	return true, ""
}
