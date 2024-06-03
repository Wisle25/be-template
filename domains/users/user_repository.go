package users

// UserRepository Interacting user-related with database
type UserRepository interface {
	AddUser(payload *RegisterUserPayload) string
	VerifyUsername(username string)
	GetUserByIdentity(identity string) (string, string)
}
