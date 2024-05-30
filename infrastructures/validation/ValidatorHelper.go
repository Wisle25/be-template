package validation

import (
	"errors"
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func fieldValue(payload interface{}, field string) interface{} {
	r := reflect.ValueOf(payload)
	f := reflect.Indirect(r).FieldByName(field)
	return f.Interface()
}

func translateError(field string, err error, trans ut.Translator) string {
	if err == nil {
		return ""
	}

	var validatorErrors validator.ValidationErrors
	errors.As(err, &validatorErrors)
	var messages []string

	for _, e := range validatorErrors {
		translated := e.Translate(trans)

		messages = append(messages, fmt.Sprintf("%s%s", field, translated))
	}

	return strings.Join(messages, ";")
}
