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

func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		return err
	}
	return nil
}
