package api

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validatePhoneNumber validator.Func = func(fieldLevel validator.FieldLevel) bool {
	isPhoneNumber := regexp.MustCompile(`^(\+234|0)\d{10}$`)
	if phoneNumber, ok := fieldLevel.Field().Interface().(string); ok {
		return isPhoneNumber.MatchString(phoneNumber)
	}
	return false
}
