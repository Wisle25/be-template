package validation

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type GoValidateUser struct {
	validator *validator.Validate
	trans     ut.Translator
}

func NewValidateUser(validator *validator.Validate, trans ut.Translator) *GoValidateUser {
	return &GoValidateUser{
		validator: validator,
		trans:     trans,
	}
}

func (v *GoValidateUser) ValidatePayload(s interface{}) {
	schema := map[string]string{
		"Username": "required,min=3,max=50,alphanum",
		"Password": "required,min=8",
		"Email":    "required,email",
	}

	for field, rule := range schema {
		value := fieldValue(s, field)

		if err := v.validator.Var(value, rule); err != nil {
			translatedErr := translateError(err, v.trans)

			panic(fmt.Errorf(translatedErr))
		}
	}
}
