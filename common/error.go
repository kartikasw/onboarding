package common

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/mattn/go-sqlite3"
)

const (
	ErrUniqueViolation = "2067"
)

var ErrRecordNotFound = sql.ErrNoRows

type Error int

const (
	ErrCredentiials Error = iota
)

func ErrorCode(err error) string {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		return fmt.Sprintf("%d", sqliteErr.ExtendedCode)
	}

	return err.Error()
}

func ErrorValidation(e error) error {
	err := e

	if errors, ok := e.(validator.ValidationErrors); ok {
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				err = fmt.Errorf("%s is required.", e.Field())
			default:
				err = fmt.Errorf("%s is invalid.", e.Field())
			}
		}
	}

	return err
}
