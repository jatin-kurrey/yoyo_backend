package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Struct(payload interface{}) map[string]string {
	err := validate.Struct(payload)
	if err == nil {
		return nil
	}

	errors := map[string]string{}
	for _, fieldError := range err.(validator.ValidationErrors) {
		field := strings.ToLower(fieldError.Field())
		switch fieldError.Tag() {
		case "required":
			errors[field] = "This field is required."
		case "email":
			errors[field] = "Enter a valid email address."
		case "min":
			errors[field] = "Value is too short."
		case "max":
			errors[field] = "Value is too long."
		case "oneof":
			errors[field] = "Value is not allowed."
		case "gte":
			errors[field] = "Value is below the minimum."
		default:
			errors[field] = "Value is invalid."
		}
	}
	return errors
}
