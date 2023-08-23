package errors

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// swagger:parameters Error
type Error struct {
	Code   string         `json:"code"`
	Errors map[string]any `json:"errors"`
}

// swagger:parameters ValidationError
type ValidationError struct {
	Code   string         `json:"code"`
	Fields map[string]any `json:"fields"`
}

func NewError(err error, code string) Error {
	e := Error{}
	e.Errors = make(map[string]any)
	e.Code = code
	switch v := err.(type) {
	case *echo.HTTPError:
		e.Errors["message"] = v.Message
	default:
		e.Errors["message"] = v.Error()
	}
	return e
}
func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "base64":
		return "Invalid base64 string"
	case "gt":
		return fmt.Sprintf("This field must be greater than %s", fe.Param())
	case "lt":
		return fmt.Sprintf("This field must be less than %s", fe.Param())
	}
	return fe.Error()
}

func NewValidatorError(err error) ValidationError {
	e := ValidationError{}
	e.Fields = make(map[string]any)
	errs := err.(validator.ValidationErrors)
	e.Code = CodeValidationErr

	for _, v := range errs {
		e.Fields[v.Field()] = fmt.Sprintf("%v", msgForTag(v))
	}
	return e
}
