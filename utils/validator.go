package utils

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"errors"
)

// STRUCT TO INDICATE VALIDATOR UTILITIES
type Validator struct {
	validate *validator.Validate
}

// CONSTRUCTOR
func NewValidator() *Validator {
	v := validator.New()

	return &Validator{validate: v}
}

// HELPER TO VALIDATE REQUIRED FIELDS AT ONCE
func (v *Validator) ValidateStruct(payload any) error {
	return v.validate.Struct(payload)
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

func FormatValidationError(err error) string {
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			return e.Field() + " is " + e.Tag()
		}
	}

	return err.Error()
}
