package validation

import (
	"net/mail"

	"github.com/go-playground/validator/v10"
)

var ValidEmail validator.Func = func(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	_, err := mail.ParseAddress(value)
	return err == nil
}
