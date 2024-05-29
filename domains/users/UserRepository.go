package users

type UserRepository interface {
	AddUser(payload *RegisterUserPayload) string
	VerifyUsername(username string)
}
