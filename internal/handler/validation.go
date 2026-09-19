package handler

import (
	"github.com/go-playground/validator/v10"
)

// formatValidationErrors converts validator errors into a user-friendly map.
func formatValidationErrors(err error) any {
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	details := make(map[string]string, len(ve))
	for _, fe := range ve {
		details[fe.Field()] = validationMessage(fe)
	}
	return details
}

// validationMessage returns a human-readable message for a validation error.
func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "min":
		return fe.Field() + " must be at least " + fe.Param()
	case "max":
		return fe.Field() + " must be at most " + fe.Param()
	case "oneof", "entity_kind", "entity_status":
		return fe.Field() + " must be one of: " + fe.Param()
	default:
		return fe.Field() + " is invalid"
	}
}
