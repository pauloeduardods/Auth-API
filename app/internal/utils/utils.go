package utils

import (
	"github.com/go-playground/validator/v10"
)

type Utils struct {
	validate *validator.Validate
}

func NewUtils(v *validator.Validate) *Utils {
	return &Utils{
		validate: v,
	}
}
