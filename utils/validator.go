package utils

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

// STRUCT TO INDICATE VALIDATOR UTILITIES
type Validator struct {
	validate *validator.Validate
}

// CONSTRUCTOR
func NewValidator() *Validator {
	return &Validator{validate: validator.New()}
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
