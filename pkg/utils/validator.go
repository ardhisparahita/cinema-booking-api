package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

func ValidationStruct(data any) error {
	return Validate.Struct(data)
}

func ValidationError(err validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for _, e := range err {
		field := strings.ToLower(e.Field())
		errs[field] = customMessage(e.Tag(), e.Param())
	}

	return errs
}

func customMessage(tag, param string) string {
	switch tag {
	case "required":
		return "this field is required"
	case "email":
		return "invalid email format"
	case "min":
		return "minimum " + param + " characters"
	case "max":
		return "maximum " + param + " characters"
	case "eqfield":
		return "does not match field " + param
	default:
		return "validation failed: " + tag
	}
}
