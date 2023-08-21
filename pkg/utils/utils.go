package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// swagger:parameters Error
type Error struct {
	Errors map[string]any `json:"errors"`
}

func NewError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]any)
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

func NewValidatorError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]any)
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		e.Errors[v.Field()] = fmt.Sprintf("%v", msgForTag(v))
	}
	return e
}

func contains[T comparable](slice []T, element T) bool {
	for _, a := range slice {
		if a == element {
			return true
		}
	}
	return false
}

// # MAY I HAVE YOUR ATTENTION, PLEASE! #
//
// DO NOT USE FOR LARGE SLICES
// Code is not optimized for performance
func IsSubSlice[T comparable](slice []T, subslice []T) bool {
	if len(slice) < len(subslice) {
		return false
	}
	for _, e := range slice {
		if !contains(subslice, e) {
			return false
		}
	}
	return true
}
