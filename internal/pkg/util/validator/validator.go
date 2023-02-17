package validator

import (
	"gitlab.miotech.com/miotech-application/backend/common-go/validator"
)

func ValidatorStruct(v interface{}) error {
	if err := validator.NewValidator().ValidateStruct(v); err != nil {
		err = validator.TranslateError(err)
		return err
	}
	return nil
}
