package repository

import "github.com/wisle25/be-template/domains/entity"

// UserRepository Interacting user-related with cache
type UserRepository interface {
	AddUser(payload *entity.RegisterUserPayload) string
	VerifyUsername(username string)
	GetUserByIdentity(identity string) (string, string)
}
