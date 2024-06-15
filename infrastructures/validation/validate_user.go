package validation

import (
	"fmt"
	"github.com/wisle25/be-template/applications/validation"
	"github.com/wisle25/be-template/domains/entity"
	"github.com/wisle25/be-template/infrastructures/services"
)

type ValidateUser struct /* implements ValidateUser */ {
	validation *services.Validation
}

func NewValidateUser(validation *services.Validation) validation.ValidateUser {
	return &ValidateUser{
		validation: validation,
	}
}

func (v *ValidateUser) ValidateRegisterPayload(payload *entity.RegisterUserPayload) {
	schema := map[string]string{
		"Username":        "required,min=3,max=50,alphanum",
		"Email":           "required,email",
		"Password":        "required,min=8",
		"ConfirmPassword": "required,min=8," + fmt.Sprintf("eq=%s", services.FieldValue(payload, "Password")),
	}

	services.Validate(payload, schema, v.validation)
}

func (v *ValidateUser) ValidateLoginPayload(payload *entity.LoginUserPayload) {
	schema := map[string]string{
		"Identity": "required,min=3,max=50",
		"Password": "required,min=8",
	}

	services.Validate(payload, schema, v.validation)
}

func (v *ValidateUser) ValidateUpdatePayload(payload *entity.UpdateUserPayload) {
	schema := map[string]string{
		"Username":        "required,min=3,max=50,alphanum",
		"Email":           "required,email",
		"Password":        "required,min=8",
		"ConfirmPassword": "required,min=8," + fmt.Sprintf("eq=%s", services.FieldValue(payload, "Password")),
	}

	services.Validate(payload, schema, v.validation)
}
