package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// GetValidationErrorMessage returns a user-friendly error message for validation errors that occur
// durinb creation of a spreadsheet.
// It uses field-specific messages for known fields and falls back to generic messages for other fields.
//
// It first checks for field-specific messages, then falls back to generic messages based on the validation tag.
func GetValidationErrorMessage(fieldErr validator.FieldError) string {
	// Field-specific messages take priority
	fieldMessages := map[string]map[string]string{
		"title": {
			"required": "Title is required",
			"max":      "Title must be 255 characters or less",
		},
		"description": {
			"required": "Description is required",
		},
		"deadline": {
			"required":    "Deadline is required",
			"time_format": "Deadline must be in format YYYY-MM-DDTHH:MM:SSZ",
		},
		"colTitles": {
			"required": "Column titles are required",
			"min":      "At least one column title is required",
		},
	}

	// Generic fallback messages for common validation tags
	tagMessages := map[string]string{
		"required":    "This field is required",
		"max":         "Must be %v characters or less",
		"min":         "Must have at least %v items",
		"time_format": "Invalid time format",
	}

	// 1. Check for field-specific message first
	if fieldMsg, ok := fieldMessages[fieldErr.Field()]; ok {
		if msg, ok := fieldMsg[fieldErr.Tag()]; ok {
			return msg
		}
	}

	// 2. Fall back to generic tag-based message
	if msg, ok := tagMessages[fieldErr.Tag()]; ok {
		// For messages that need the validation parameter (like max:255)
		if fieldErr.Param() != "" {
			return fmt.Sprintf(msg, fieldErr.Param())
		}
		return msg
	}

	// 3. Final fallback
	return fmt.Sprintf("Validation failed on field '%s'", fieldErr.Field())
}
