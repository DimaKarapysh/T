package bootstrap

import (
	"T/internal/validator"
)

func InitValidator() validator.Validator {
	return validator.NewValidator()
}
