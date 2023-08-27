package utils

import (
	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()
}

func Validate(s interface{}) *ApiError {
	err := validate.Struct(s)
	if err != nil {
		return NewApiError(400, "Invalid request body", err.Error())
	}
	return nil
}
