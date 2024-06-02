package validation

// ValidateUser interface defines methods for validating user-related payloads.
type ValidateUser interface {
	ValidateRegisterPayload(s interface{})
	ValidateLoginPayload(s interface{})
}
