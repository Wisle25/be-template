package users

type UserRepository interface {
	AddUser(payload *RegisterUserPayload) *RegisterUserResponse
	VerifyUsername(username string)
}
