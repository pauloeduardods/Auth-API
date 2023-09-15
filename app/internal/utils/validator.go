package utils

func (v *Utils) Validate(s interface{}) error {
	err := v.validate.Struct(s)
	if err != nil {
		return err
	}
	return nil
}
