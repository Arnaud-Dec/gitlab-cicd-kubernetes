package server

import (
	"github.com/go-playground/validator/v10"
)

// EchoValidator adapts go-playground/validator to echo.Validator.
type EchoValidator struct {
	validator *validator.Validate
}

func NewValidator() *EchoValidator {
	return &EchoValidator{validator: validator.New()}
}

func (v *EchoValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
