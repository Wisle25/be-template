package use_case

import (
	"github.com/wisle25/be-template/applications/cache"
	"github.com/wisle25/be-template/applications/processing"
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
	fileProcessing processing.FileProcessing
	passwordHash   security.PasswordHash
	validator      validation.ValidateUser
	config         *commons.Config
	token          security.Token
	cache          cache.Cache
}

func NewUserUseCase(
	userRepository repository.UserRepository,
	fileProcessing processing.FileProcessing,
	passwordHash security.PasswordHash,
	validator validation.ValidateUser,
	config *commons.Config,
	token security.Token,
	cache cache.Cache,
) *UserUseCase {
	return &UserUseCase{
		userRepository: userRepository,
		fileProcessing: fileProcessing,
		passwordHash:   passwordHash,
		validator:      validator,
		config:         config,
		token:          token,
		cache:          cache,
	}
}

// ExecuteAdd Handling user registration.
// Should raise panic if violates username/email uniqueness.
// Returning registered user's ID.
func (uc *UserUseCase) ExecuteAdd(payload *entity.RegisterUserPayload) string {
	uc.validator.ValidateRegisterPayload(payload)

	payload.Password = uc.passwordHash.Hash(payload.Password)

	return uc.userRepository.AddUser(payload)
}

// ExecuteLogin Handling user login. Returning user's token for authentication/authorization later.
// Should raise panic if user is not existed
// Returned tokens must be added to the HTTP cookie
func (uc *UserUseCase) ExecuteLogin(payload *entity.LoginUserPayload) (*entity.TokenDetail, *entity.TokenDetail) {
	uc.validator.ValidateLoginPayload(payload)

	// Get user information from database then compare password
	userId, encryptedPassword := uc.userRepository.GetUserForLogin(payload.Identity)
	uc.passwordHash.Compare(payload.Password, encryptedPassword)

	// Create token
	accessTokenDetail := uc.token.CreateToken(userId, uc.config.AccessTokenExpiresIn, uc.config.AccessTokenPrivateKey)
	refreshTokenDetail := uc.token.CreateToken(userId, uc.config.RefreshTokenExpiresIn, uc.config.RefreshTokenPrivateKey)

	// Add tokens to the cache
	now := time.Now()
	uc.cache.SetCache(accessTokenDetail.TokenId, userId, time.Unix(accessTokenDetail.ExpiresIn, 0).Sub(now))
	uc.cache.SetCache(refreshTokenDetail.TokenId, userId, time.Unix(refreshTokenDetail.ExpiresIn, 0).Sub(now))

	// Returned token should be added to HTTP Cookie
	return accessTokenDetail, refreshTokenDetail
}

// ExecuteRefreshToken handles refreshing the access token using the provided refresh token.
// Should raise panic if refresh token is invalid
// Returned new access token should be added to HTTP Cookie
func (uc *UserUseCase) ExecuteRefreshToken(currentRefreshToken string) *entity.TokenDetail {
	// Verify token from JWT itself and from cache
	tokenClaims := uc.token.ValidateToken(currentRefreshToken, uc.config.RefreshTokenPublicKey)
	userId := uc.cache.GetCache(tokenClaims.TokenId).(string)

	// Re-create access token and re-insert to the cache
	now := time.Now()
	accessTokenDetail := uc.token.CreateToken(userId, uc.config.AccessTokenExpiresIn, uc.config.AccessTokenPrivateKey)
	uc.cache.SetCache(accessTokenDetail.TokenId, userId, time.Unix(accessTokenDetail.ExpiresIn, 0).Sub(now))

	// Returned token should be added to HTTP Cookie
	return accessTokenDetail
}

// ExecuteLogout handles user logout by removing the tokens from the cache.
// Don't forget to remove the tokens from cookies too in infrastructure layer
func (uc *UserUseCase) ExecuteLogout(refreshToken string, accessTokenId string) {
	// Verify
	refreshTokenClaims := uc.token.ValidateToken(refreshToken, uc.config.RefreshTokenPublicKey)

	// Remove from cache
	uc.cache.DeleteCache(refreshTokenClaims.TokenId)
	uc.cache.DeleteCache(accessTokenId)
}

// ExecuteGuard verifies the access token and retrieves the associated user from the cache.
// This is used as a guard middleware for JWT authentication.
// Returning userId from token's cache
func (uc *UserUseCase) ExecuteGuard(accessToken string) (interface{}, *entity.TokenDetail) {
	accessTokenDetail := uc.token.ValidateToken(accessToken, uc.config.AccessTokenPublicKey)

	return uc.cache.GetCache(accessTokenDetail.TokenId), accessTokenDetail
}

// ExecuteGetUserById simply returns specified user information by ID
func (uc *UserUseCase) ExecuteGetUserById(userId string) *entity.User {
	return uc.userRepository.GetUserById(userId)
}

// ExecuteUpdateUserById Updating user information and now user can set their new password and upload an avatar.
func (uc *UserUseCase) ExecuteUpdateUserById(userId string, payload *entity.UpdateUserPayload) {
	uc.validator.ValidateUpdatePayload(payload)
	
	// Hash password
	payload.Password = uc.passwordHash.Hash(payload.Password)

	// Handling avatar file
	newAvatarLink := uc.fileProcessing.UploadFile(payload.Avatar)
	oldAvatarLink := uc.userRepository.UpdateUserById(userId, payload, newAvatarLink)

	if oldAvatarLink != "" {
		uc.fileProcessing.RemoveFile(oldAvatarLink)
	}
}
