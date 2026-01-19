package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(obj interface{}) error {
	err := validate.Struct(obj)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors.New("validation error")
	}

	if len(validationErrors) > 0 {
		return FormatValidationError(validationErrors[0])
	}

	return errors.New("validation error")
}

func ValidateStructAll(obj interface{}) []string {
	err := validate.Struct(obj)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return []string{"validation error"}
	}

	var errors []string
	for _, validationError := range validationErrors {
		errors = append(errors, FormatValidationError(validationError).Error())
	}

	return errors
}

func FormatValidationError(validationError validator.FieldError) error {
	field := strings.ToLower(validationError.StructField())

	switch validationError.Tag() {
	case "required":
		return fmt.Errorf("%s is required", field)

	case "min":
		return fmt.Errorf("%s must be at least %s characters", field, validationError.Param())

	case "max":
		return fmt.Errorf("%s must be at most %s characters", field, validationError.Param())

	case "len":
		return fmt.Errorf("%s must be exactly %s characters", field, validationError.Param())

	case "email":
		return fmt.Errorf("%s must be a valid email address", field)

	case "numeric":
		return fmt.Errorf("%s must contain only numbers", field)

	case "alpha":
		return fmt.Errorf("%s must contain only letters", field)

	case "alphanum":
		return fmt.Errorf("%s must contain only letters and numbers", field)

	case "oneof":
		return fmt.Errorf("%s must be one of [%s]", field, validationError.Param())

	case "cpf":
		return fmt.Errorf("%s must be a valid CPF", field)

	case "phone_br":
		return fmt.Errorf("%s must be a valid Brazilian phone number", field)

	case "gt":
		return fmt.Errorf("%s must be greater than %s", field, validationError.Param())

	case "gte":
		return fmt.Errorf("%s must be greater than or equal to %s", field, validationError.Param())

	case "lt":
		return fmt.Errorf("%s must be less than %s", field, validationError.Param())

	case "lte":
		return fmt.Errorf("%s must be less than or equal to %s", field, validationError.Param())

	default:
		return fmt.Errorf("%s is invalid", field)
	}
}

func GetValidator() *validator.Validate {
	return validate
}
