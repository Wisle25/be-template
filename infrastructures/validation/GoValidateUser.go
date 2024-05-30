package validation

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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

func (v *GoValidateUser) ValidatePayload(s interface{}) {
	schema := map[string]string{
		"Username":        "required,min=3,max=50,alphanum",
		"Password":        "required,min=8",
		"Email":           "required,email",
		"ConfirmPassword": "required,min=8," + fmt.Sprintf("eq=%s", fieldValue(s, "Password")),
	}

	for field, rule := range schema {
		value := fieldValue(s, field)

		if err := v.validator.Var(value, rule); err != nil {
			translatedErr := translateError(field, err, v.trans)

			if field == "ConfirmPassword" {
				translatedErr = "Confirm Password doesn't match!"
			}

			panic(fiber.NewError(fiber.StatusBadRequest, translatedErr))
		}
	}
}
