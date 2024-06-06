package use_case_test

import (
	"github.com/wisle25/be-template/applications/use_case"
	"github.com/wisle25/be-template/domains/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wisle25/be-template/commons"
)

// Mocks for the dependencies

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) AddUser(payload *entity.RegisterUserPayload) string {
	args := m.Called(payload)
	return args.String(0)
}

func (m *MockUserRepository) VerifyUsername(username string) {
	m.Called(username)
}

func (m *MockUserRepository) GetUserByIdentity(identity string) (string, string) {
	args := m.Called(identity)
	return args.String(0), args.String(1)
}

type MockPasswordHash struct {
	mock.Mock
}

func (m *MockPasswordHash) Hash(password string) string {
	args := m.Called(password)
	return args.String(0)
}

func (m *MockPasswordHash) Compare(password string, hashedPassword string) {
	m.Called(password, hashedPassword)
}

type MockValidateUser struct {
	mock.Mock
}

func (m *MockValidateUser) ValidateRegisterPayload(s interface{}) {
	m.Called(s)
}

func (m *MockValidateUser) ValidateLoginPayload(s interface{}) {
	m.Called(s)
}

type MockToken struct {
	mock.Mock
}

func (m *MockToken) CreateToken(userID string, ttl time.Duration, privateKey string) *entity.TokenDetail {
	args := m.Called(userID, ttl, privateKey)
	return args.Get(0).(*entity.TokenDetail)
}

func (m *MockToken) ValidateToken(token string, publicKey string) *entity.TokenDetail {
	args := m.Called(token, publicKey)
	return args.Get(0).(*entity.TokenDetail)
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) SetCache(key string, value interface{}, expiration time.Duration) {
	m.Called(key, value, expiration)
}

func (m *MockCache) GetCache(key string) interface{} {
	args := m.Called(key)
	return args.Get(0)
}

func (m *MockCache) DeleteCache(key string) {
	m.Called(key)
}

func TestUserUseCase_ExecuteAdd(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockPasswordHash := new(MockPasswordHash)
	mockValidator := new(MockValidateUser)
	mockConfig := &commons.Config{}
	mockToken := new(MockToken)
	mockCache := new(MockCache)

	userUseCase := use_case.NewUserUseCase(mockUserRepo, mockPasswordHash, mockValidator, mockConfig, mockToken, mockCache)

	payload := &entity.RegisterUserPayload{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}

	mockValidator.On("ValidateRegisterPayload", payload).Return(nil)
	mockUserRepo.On("VerifyUsername", payload.Username).Return(nil)
	mockPasswordHash.On("Hash", payload.Password).Return("hashedpassword")
	mockUserRepo.On("AddUser", payload).Return("userid123")

	// Action
	userId := userUseCase.ExecuteAdd(payload)

	// Assert
	assert.Equal(t, "userid123", userId)

	mockValidator.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockPasswordHash.AssertExpectations(t)
}

func TestUserUseCase_ExecuteLogin(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockPasswordHash := new(MockPasswordHash)
	mockValidator := new(MockValidateUser)
	mockConfig := &commons.Config{
		AccessTokenExpiresIn:  time.Hour,
		RefreshTokenExpiresIn: time.Hour * 24,
	}
	mockToken := new(MockToken)
	mockCache := new(MockCache)

	userUseCase := use_case.NewUserUseCase(mockUserRepo, mockPasswordHash, mockValidator, mockConfig, mockToken, mockCache)

	payload := &entity.LoginUserPayload{
		Identity: "testuser",
		Password: "password123",
	}

	user := &entity.User{
		Id:       "userid123",
		Username: "testuser",
		Email:    "test@example.com",
	}

	accessTokenDetail := &entity.TokenDetail{
		TokenID:   "access_token_id",
		ExpiresIn: time.Now().Add(time.Hour).Unix(),
		UserID:    "userid123",
		Token:     "access_token",
	}

	refreshTokenDetail := &entity.TokenDetail{
		TokenID:   "refresh_token_id",
		ExpiresIn: time.Now().Add(time.Hour * 24).Unix(),
		UserID:    "userid123",
		Token:     "refresh_token",
	}

	mockValidator.On("ValidateLoginPayload", payload).Return(nil)
	mockUserRepo.On("GetUserByIdentity", payload.Identity).Return(user.Id, "hashedpassword")
	mockPasswordHash.On("Compare", payload.Password, "hashedpassword").Return(nil)
	mockToken.On("CreateToken", user.Id, mockConfig.AccessTokenExpiresIn, mockConfig.AccessTokenPrivateKey).Return(accessTokenDetail)
	mockToken.On("CreateToken", user.Id, mockConfig.RefreshTokenExpiresIn, mockConfig.RefreshTokenPrivateKey).Return(refreshTokenDetail)
	mockCache.On("SetCache", accessTokenDetail.TokenID, user.Id, mock.Anything).Return(nil)
	mockCache.On("SetCache", refreshTokenDetail.TokenID, user.Id, mock.Anything).Return(nil)

	// Action
	accessToken, refreshToken := userUseCase.ExecuteLogin(payload)

	// Assert
	assert.Equal(t, accessTokenDetail, accessToken)
	assert.Equal(t, refreshTokenDetail, refreshToken)

	mockValidator.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockPasswordHash.AssertExpectations(t)
	mockToken.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUserUseCase_ExecuteRefreshToken(t *testing.T) {
	// Arrange
	mockToken := new(MockToken)
	mockCache := new(MockCache)
	mockConfig := &commons.Config{
		AccessTokenExpiresIn:  time.Hour,
		AccessTokenPrivateKey: "any",
		RefreshTokenPublicKey: "any",
	}

	userUseCase := use_case.NewUserUseCase(
		&MockUserRepository{},
		&MockPasswordHash{},
		&MockValidateUser{},
		mockConfig,
		mockToken,
		mockCache,
	)

	refreshTokenCookie := "refresh_token123"

	accessTokenDetail := &entity.TokenDetail{
		TokenID:   "access_token_id",
		ExpiresIn: time.Now().Add(time.Hour).Unix(),
		UserID:    "userid123",
		Token:     "access_token",
	}
	refreshTokenDetail := &entity.TokenDetail{
		TokenID:   "refresh_token_id",
		ExpiresIn: time.Now().Add(time.Hour * 24).Unix(),
		UserID:    "userid123",
		Token:     "refresh_token",
	}

	mockToken.On("ValidateToken", refreshTokenCookie, mockConfig.RefreshTokenPublicKey).Return(refreshTokenDetail)
	mockCache.On("GetCache", refreshTokenDetail.TokenID).Return(refreshTokenDetail.UserID)
	mockToken.On("CreateToken", refreshTokenDetail.UserID, mockConfig.AccessTokenExpiresIn, mockConfig.AccessTokenPrivateKey).Return(accessTokenDetail)
	mockCache.On("SetCache", accessTokenDetail.TokenID, refreshTokenDetail.UserID, mock.Anything).Return(nil)

	// Action
	accessTokenResponse := userUseCase.ExecuteRefreshToken(refreshTokenCookie)

	// Assert
	assert.Equal(t, accessTokenDetail, accessTokenResponse)
	mockToken.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUserUseCase_ExecuteLogout(t *testing.T) {
	// Arrange
	refreshTokenCookie := "refresh_token123"
	accessTokenId := "access_token123"

	refreshTokenDetail := &entity.TokenDetail{
		TokenID: "refresh_token123",
	}

	mockConfig := &commons.Config{
		RefreshTokenPublicKey: "any",
	}
	mockToken := new(MockToken)
	mockCache := new(MockCache)

	mockToken.On("ValidateToken", refreshTokenCookie, mockConfig.RefreshTokenPublicKey).Return(refreshTokenDetail)
	mockCache.On("DeleteCache", refreshTokenDetail.TokenID).Return(nil).Once()
	mockCache.On("DeleteCache", accessTokenId).Return(nil).Once()

	userUseCase := use_case.NewUserUseCase(
		&MockUserRepository{},
		&MockPasswordHash{},
		&MockValidateUser{},
		mockConfig,
		mockToken,
		mockCache,
	)

	// Action
	userUseCase.ExecuteLogout(refreshTokenCookie, accessTokenId)

	// Assert
	mockToken.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestUseUseCase_ExecuteGuard(t *testing.T) {
	// Arrange
	accessToken := "access_token123"
	accessTokenDetail := &entity.TokenDetail{
		TokenID: "access_token123",
		UserID:  "userid123",
	}

	mockToken := new(MockToken)
	mockCache := new(MockCache)
	mockConfig := &commons.Config{
		RefreshTokenPublicKey: "any",
	}

	mockToken.On("ValidateToken", accessToken, mockConfig.AccessTokenPublicKey).Return(accessTokenDetail)
	mockCache.On("GetCache", accessTokenDetail.TokenID).Return(accessTokenDetail.UserID)

	userUseCase := use_case.NewUserUseCase(
		&MockUserRepository{},
		&MockPasswordHash{},
		&MockValidateUser{},
		mockConfig,
		mockToken,
		mockCache,
	)

	// Action
	userIdCache, tokenDetail := userUseCase.ExecuteGuard(accessToken)

	assert.NotNil(t, userIdCache)
	assert.Equal(t, accessTokenDetail, tokenDetail)
	mockToken.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}