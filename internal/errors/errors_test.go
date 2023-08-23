package errors

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	err := NewError(errors.New("unknown"), CodeUnknownErr)

	assert.Equal(t, CodeUnknownErr, err.Code)
	assert.Equal(t, "unknown", err.Errors["message"])
}

func TestNewValidatorError(t *testing.T) {
	var e validator.ValidationErrors = []validator.FieldError{}

	err := NewValidatorError(e)

	assert.Equal(t, CodeValidationErr, err.Code)
}
