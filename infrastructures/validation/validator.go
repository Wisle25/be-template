package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
)

func NewValidator() *validator.Validate {
	validate := validator.New()

	return validate
}

func NewValidatorTranslator(validate *validator.Validate) ut.Translator {
	// Using translator to translate error message
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = en2.RegisterDefaultTranslations(validate, trans)

	// Overriding "eq" translation
	err := validate.RegisterTranslation("eq", trans, func(ut ut.Translator) error {
		return ut.Add("eq", "{0} doesn't match!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("eq", fe.Field())

		return t
	})
	if err != nil {
		panic(err)
	}

	return trans
}
