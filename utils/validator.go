package utils

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"regexp"
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
func (v *Validator) ValidateStruct(payload any) (string, error) {
	var err error
	var message string

	err = v.validate.Struct(payload)

	if err != nil {
		message = FormatValidationError(err, payload)
	}

	return message, err
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

func FormatValidationError(err error, payload any) string {
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			if e.Tag() == "required" {
				return e.Field() + " is required"
			}

			if e.Tag() == "ne" {
				return e.Field() + " cannot be null"
			}
		}
	}

	return err.Error()
}
