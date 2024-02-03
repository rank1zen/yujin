package main

import "github.com/go-playground/validator/v10"

type CValidator struct {
	v *validator.Validate
}

func NewValidator() *CValidator {
	return &CValidator{v: validator.New()}
}

func (cv *CValidator) Validate(i interface{}) error {
	return cv.v.Struct(i)
}
