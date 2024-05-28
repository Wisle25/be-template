package validation

type ValidateUser interface {
	ValidatePayload(s interface{})
}
