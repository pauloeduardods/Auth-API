package utils

import (
	"github.com/go-playground/validator/v10"
)

type Utils struct {
	validate *validator.Validate
}

type UtilsOptions struct {
	Validate *validator.Validate
}

func NewUtils(opts UtilsOptions) *Utils {
	return &Utils{
		validate: opts.Validate,
	}
}
