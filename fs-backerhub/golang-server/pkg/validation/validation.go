package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ExtractValidationError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		fe := errs[0]

		switch fe.Tag() {
		case "required":
			return fmt.Sprintf("%s is required", fe.Field())
		case "email":
			return fmt.Sprintf("%s must be a valid email", fe.Field())
		case "min":
			return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
		case "max":
			return fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
		default:
			return fmt.Sprintf("%s is invalid", fe.Field())
		}
	}

	return "Invalid request payload"
}
