package validation

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/wisle25/be-template/applications/validation"
)

type GoValidateUser struct {
	validator *validator.Validate
	trans     ut.Translator
}

func NewValidateUser(validator *validator.Validate, trans ut.Translator) validation.ValidateUser {
	return &GoValidateUser{
		validator: validator,
		trans:     trans,
	}
}

func (v *GoValidateUser) ValidateRegisterPayload(s interface{}) {
	schema := map[string]string{
		"Username":        "required,min=3,max=50,alphanum",
		"Password":        "required,min=8",
		"Email":           "required,email",
		"ConfirmPassword": "required,min=8," + fmt.Sprintf("eq=%s", fieldValue(s, "Password")),
	}

	validate(s, schema, v)
}

func (v *GoValidateUser) ValidateLoginPayload(s interface{}) {
	schema := map[string]string{
		"Identity": "required,min=3,max=50",
		"Password": "required,min=8",
	}

	validate(s, schema, v)
}
