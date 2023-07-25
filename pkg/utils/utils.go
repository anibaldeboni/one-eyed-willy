package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// swagger:parameters Error
type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

func NewError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	switch v := err.(type) {
	case *echo.HTTPError:
		e.Errors["body"] = v.Message
	default:
		e.Errors["body"] = v.Error()
	}
	return e
}

func NewValidatorError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		e.Errors[v.Field()] = fmt.Sprintf("%v", v.Tag())
	}
	return e
}

func AccessForbidden() Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = "access forbidden"
	return e
}

func NotFound() Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = "resource not found"
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
