package validatorUtil

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

type Options struct {
	Validate *validator.Validate
}

func New(opts Options) *Validator {
	return &Validator{
		validate: opts.Validate,
	}
}

func (v *Validator) Validate(s interface{}) error {
	err := v.validate.Struct(s)
	if err != nil {
		return err
	}
	return nil
}
