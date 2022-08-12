package validator

import (
	"mio/pkg/validator"
)

func ValidatorStruct(v interface{}) error {
	if err := validator.NewValidator().ValidateStruct(v); err != nil {
		err = validator.TranslateError(err)
		return err
	}
	return nil
}
