package api

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterJSONTagNameFormatter sets a custom tag name function for JSON binding validation.
//
// This function configures the validator to use the `json` struct tag as the field name
// when returning validation errors. Specifically, it affects the output of FieldError.Field(),
// so that it returns the JSON field name (e.g., `title`) instead of the Go struct field name (e.g., `Title`).
//
// This improves the clarity of validation error messages when working with JSON APIs.
func RegisterJSONTagNameFormatter() string {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
	return ""
}
