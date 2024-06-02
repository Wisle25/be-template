package use_case

import (
	"github.com/wisle25/be-template/applications/database"
	"github.com/wisle25/be-template/applications/security"
	"github.com/wisle25/be-template/applications/validation"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/domains/tokens"
	"github.com/wisle25/be-template/domains/users"
	"time"
)

// UserUseCase handles the business logic for user operations.
type UserUseCase struct {
	userRepository users.UserRepository
	passwordHash   security.PasswordHash
	validator      validation.ValidateUser
	config         *commons.Config
	token          security.Token
	cache          database.Cache
}

func NewUserUseCase(
	userRepository users.UserRepository,
	passwordHash security.PasswordHash,
	validator validation.ValidateUser,
	config *commons.Config,
	token security.Token,
	cache database.Cache,
) *UserUseCase {
	return &UserUseCase{
		userRepository: userRepository,
		passwordHash:   passwordHash,
		validator:      validator,
		config:         config,
		token:          token,
		cache:          cache,
	}
}

// ExecuteAdd Handling user registration. Returning registered user's ID
func (uc *UserUseCase) ExecuteAdd(payload *users.RegisterUserPayload) string {
	uc.validator.ValidateRegisterPayload(payload)

	uc.userRepository.VerifyUsername(payload.Username)
	payload.Password = uc.passwordHash.Hash(payload.Password)

	return uc.userRepository.AddUser(payload)
}

// ExecuteLogin Handling user login. Returning user's token for authentication/authorization later.
// Returned token must be added to the cookie
func (uc *UserUseCase) ExecuteLogin(payload *users.LoginUserPayload) (*tokens.TokenDetail, *tokens.TokenDetail) {
	uc.validator.ValidateLoginPayload(payload)

	user, encryptedPassword := uc.userRepository.GetUserByIdentity(payload.Identity)
	uc.passwordHash.Compare(payload.Password, encryptedPassword)

	accessTokenDetail := uc.token.CreateToken(user.Id, uc.config.AccessTokenExpiresIn, uc.config.AccessTokenPrivateKey)
	refreshTokenDetail := uc.token.CreateToken(user.Id, uc.config.RefreshTokenExpiresIn, uc.config.RefreshTokenPrivateKey)

	now := time.Now()
	uc.cache.SetCache(accessTokenDetail.TokenID, user.Id, time.Unix(accessTokenDetail.ExpiresIn, 0).Sub(now))
	uc.cache.SetCache(refreshTokenDetail.TokenID, user.Id, time.Unix(refreshTokenDetail.ExpiresIn, 0).Sub(now))

	return accessTokenDetail, refreshTokenDetail
}

func (uc *UserUseCase) ExecuteRefreshToken(payload string) *tokens.TokenDetail {
	// Verify
	tokenClaims := uc.token.ValidateToken(payload, uc.config.RefreshTokenPublicKey)
	userId := uc.cache.GetCache(tokenClaims.TokenID).(string)

	// Re-create access token
	now := time.Now()

	accessTokenDetail := uc.token.CreateToken(userId, uc.config.AccessTokenExpiresIn, uc.config.AccessTokenPrivateKey)
	uc.cache.SetCache(accessTokenDetail.TokenID, userId, time.Unix(accessTokenDetail.ExpiresIn, 0).Sub(now))

	return accessTokenDetail
}
