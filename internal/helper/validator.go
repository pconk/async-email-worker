package helper

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// FormatValidationError mengubah error validator mentah jadi map yang rapi
func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	// Cek apakah error memang dari validator
	if castedObject, ok := err.(validator.ValidationErrors); ok {
		for _, err := range castedObject {
			// Contoh: field "Name" dengan tag "required" jadi "Name: is required"
			switch err.Tag() {
			case "required":
				errors[err.Field()] = "is required"
			case "min":
				errors[err.Field()] = fmt.Sprintf("must be at least %s characters", err.Param())
			case "gt":
				errors[err.Field()] = fmt.Sprintf("must be greater than %s", err.Param())
			default:
				errors[err.Field()] = "is invalid"
			}
		}
	}
	return errors
}
