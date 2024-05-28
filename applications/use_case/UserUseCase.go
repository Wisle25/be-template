package use_case

import (
	"github.com/wisle25/be-template/applications/security"
	"github.com/wisle25/be-template/applications/validation"
	"github.com/wisle25/be-template/domains/users"
)

type UserUseCase struct {
	userRepository users.UserRepository
	passwordHash   security.PasswordHash
	validator      validation.ValidateUser
}

func NewAddUserUseCase(
	userRepository users.UserRepository,
	passwordHash security.PasswordHash,
	validator validation.ValidateUser,
) *UserUseCase {
	return &UserUseCase{
		userRepository: userRepository,
		passwordHash:   passwordHash,
		validator:      validator,
	}
}

func (uc *UserUseCase) Execute(payload *users.RegisterUserPayload) *users.RegisterUserResponse {
	uc.validator.ValidatePayload(payload)

	uc.userRepository.VerifyUsername(payload.Username)
	payload.Password = uc.passwordHash.Hash(payload.Password)

	return uc.userRepository.AddUser(payload)
}
