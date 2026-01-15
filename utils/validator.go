package utils

import (
	"errors"
	"regexp"
)

// STRUCT TO INDICATE VALIDATOR UTILITIES
type Validator struct{}

// HELPER TO VALIDATE REQUIRED FIELDS AT ONCE
func (v *Validator) ValidateRequiredFields(fields map[string]any) error {
	for name, value := range fields {
		if ok, message := v.validateRequired(value, name); !ok {
			return errors.New(message)
		}
	}

	return nil
}

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

// VALIDATE PHONE NUMBER FORMAT
func (v *Validator) ValidatePhone(phone string) bool {
	if phone == "" {
		return true
	}

	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

	return phoneRegex.MatchString(phone)
}

// REQUIRED FIELD VALIDATION
func (v *Validator) validateRequired(value any, fieldName string) (bool, string) {
	if value == nil {
		return false, fieldName + " is required"
	}

	switch v := value.(type) {
	case string:
		if v == "" {
			return false, fieldName + " is required"
		}
	}

	return true, ""
}
