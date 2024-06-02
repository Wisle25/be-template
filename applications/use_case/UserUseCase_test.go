package use_case_test

import (
	"github.com/wisle25/be-template/applications/use_case"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/domains/tokens"
	"github.com/wisle25/be-template/domains/users"
)

// Mocks for the dependencies

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) AddUser(payload *users.RegisterUserPayload) string {
	args := m.Called(payload)
	return args.String(0)
}

func (m *MockUserRepository) VerifyUsername(username string) {
	m.Called(username)
}

func (m *MockUserRepository) GetUserByIdentity(identity string) (*users.User, string) {
	args := m.Called(identity)
	return args.Get(0).(*users.User), args.String(1)
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

func (m *MockToken) CreateToken(userID string, ttl time.Duration, privateKey string) *tokens.TokenDetail {
	args := m.Called(userID, ttl, privateKey)
	return args.Get(0).(*tokens.TokenDetail)
}

func (m *MockToken) ValidateToken(token string, publicKey string) *tokens.TokenDetail {
	args := m.Called(token, publicKey)
	return args.Get(0).(*tokens.TokenDetail)
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

func TestUserUseCase_ExecuteAdd(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockPasswordHash := new(MockPasswordHash)
	mockValidator := new(MockValidateUser)
	mockConfig := &commons.Config{}
	mockToken := new(MockToken)
	mockCache := new(MockCache)

	userUseCase := use_case.NewUserUseCase(mockUserRepo, mockPasswordHash, mockValidator, mockConfig, mockToken, mockCache)

	payload := &users.RegisterUserPayload{
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

	payload := &users.LoginUserPayload{
		Identity: "testuser",
		Password: "password123",
	}

	user := &users.User{
		Id:       "userid123",
		Username: "testuser",
		Email:    "test@example.com",
	}

	accessTokenDetail := &tokens.TokenDetail{
		TokenID:   "access_token_id",
		ExpiresIn: time.Now().Add(time.Hour).Unix(),
		UserID:    "userid123",
		Token:     "access_token",
	}

	refreshTokenDetail := &tokens.TokenDetail{
		TokenID:   "refresh_token_id",
		ExpiresIn: time.Now().Add(time.Hour * 24).Unix(),
		UserID:    "userid123",
		Token:     "refresh_token",
	}

	mockValidator.On("ValidateLoginPayload", payload).Return(nil)
	mockUserRepo.On("GetUserByIdentity", payload.Identity).Return(user, "hashedpassword")
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
