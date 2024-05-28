package validator

type ValidateUser interface {
	ValidatePayload(s interface{})
}