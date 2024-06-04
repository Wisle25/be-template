package use_case

import (
	"github.com/wisle25/be-template/applications/cache"
	"github.com/wisle25/be-template/applications/security"
	"github.com/wisle25/be-template/applications/validation"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/domains/entity"
	"github.com/wisle25/be-template/domains/repository"
	"time"
)

// UserUseCase handles the business logic for user operations.
type UserUseCase struct {
	userRepository repository.UserRepository
	passwordHash   security.PasswordHash
	validator      validation.ValidateUser
	config         *commons.Config
	token          security.Token
	cache          cache.Cache
}

func NewUserUseCase(
	userRepository repository.UserRepository,
	passwordHash security.PasswordHash,
	validator validation.ValidateUser,
	config *commons.Config,
	token security.Token,
	cache cache.Cache,
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
func (uc *UserUseCase) ExecuteAdd(payload *entity.RegisterUserPayload) string {
	uc.validator.ValidateRegisterPayload(payload)

	uc.userRepository.VerifyUsername(payload.Username)
	payload.Password = uc.passwordHash.Hash(payload.Password)

	return uc.userRepository.AddUser(payload)
}

// ExecuteLogin Handling user login. Returning user's token for authentication/authorization later.
// Returned token must be added to the cookie
func (uc *UserUseCase) ExecuteLogin(payload *entity.LoginUserPayload) (*entity.TokenDetail, *entity.TokenDetail) {
	uc.validator.ValidateLoginPayload(payload)

	userId, encryptedPassword := uc.userRepository.GetUserByIdentity(payload.Identity)
	uc.passwordHash.Compare(payload.Password, encryptedPassword)

	accessTokenDetail := uc.token.CreateToken(userId, uc.config.AccessTokenExpiresIn, uc.config.AccessTokenPrivateKey)
	refreshTokenDetail := uc.token.CreateToken(userId, uc.config.RefreshTokenExpiresIn, uc.config.RefreshTokenPrivateKey)

	now := time.Now()
	uc.cache.SetCache(accessTokenDetail.TokenID, userId, time.Unix(accessTokenDetail.ExpiresIn, 0).Sub(now))
	uc.cache.SetCache(refreshTokenDetail.TokenID, userId, time.Unix(refreshTokenDetail.ExpiresIn, 0).Sub(now))

	return accessTokenDetail, refreshTokenDetail
}

func (uc *UserUseCase) ExecuteRefreshToken(currentRefreshToken string) *entity.TokenDetail {
	// Verify
	tokenClaims := uc.token.ValidateToken(currentRefreshToken, uc.config.RefreshTokenPublicKey)
	userId := uc.cache.GetCache(tokenClaims.TokenID).(string)

	// Re-create access token
	now := time.Now()

	accessTokenDetail := uc.token.CreateToken(userId, uc.config.AccessTokenExpiresIn, uc.config.AccessTokenPrivateKey)
	uc.cache.SetCache(accessTokenDetail.TokenID, userId, time.Unix(accessTokenDetail.ExpiresIn, 0).Sub(now))

	return accessTokenDetail
}

func (uc *UserUseCase) ExecuteLogout(refreshToken string, accessTokenId string) {
	// Verify
	refreshTokenClaims := uc.token.ValidateToken(refreshToken, uc.config.RefreshTokenPublicKey)

	// Remove from cache
	uc.cache.DeleteCache(refreshTokenClaims.TokenID)
	uc.cache.DeleteCache(accessTokenId)
}

func (uc *UserUseCase) ExecuteGuard(accessToken string) (interface{}, *entity.TokenDetail) {
	accessTokenDetail := uc.token.ValidateToken(accessToken, uc.config.AccessTokenPublicKey)

	return uc.cache.GetCache(accessTokenDetail.TokenID), accessTokenDetail
}
