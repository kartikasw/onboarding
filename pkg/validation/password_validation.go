package validation

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type PasswordValidationError struct {
	Reasons []string
}

func (e PasswordValidationError) Error() string {
	return "Password is invalid: " + strings.Join(e.Reasons, ", ")
}

var ValidPassword validator.Func = func(fl validator.FieldLevel) bool {
	pass, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(pass) >= 8 {
		hasMinLen = true
	}

	for _, c := range pass {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
