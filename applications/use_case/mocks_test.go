package use_case_test

import (
	"github.com/stretchr/testify/mock"
	"github.com/wisle25/be-template/applications/file_handling"
	"github.com/wisle25/be-template/domains/entity"
	"time"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) RegisterUser(payload *entity.RegisterUserPayload) {
	m.Called(payload)
}

func (m *MockUserRepository) GetAllUsers() []entity.User {
	args := m.Called()

	return args.Get(0).([]entity.User)
}

func (m *MockUserRepository) GetUserForLogin(identity string) (*entity.User, string) {
	args := m.Called(identity)

	return args.Get(0).(*entity.User), args.String(1)
}

func (m *MockUserRepository) GetUserById(id string) *entity.User {
	args := m.Called(id)

	return args.Get(0).(*entity.User)
}

func (m *MockUserRepository) UpdateUserById(id string, payload *entity.UpdateUserPayload, newAvatarLink string) string {
	args := m.Called(id, payload, newAvatarLink)

	return args.String(0)
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

func (m *MockValidateUser) ValidateRegisterPayload(payload *entity.RegisterUserPayload) {
	m.Called(payload)
}

func (m *MockValidateUser) ValidateLoginPayload(payload *entity.LoginUserPayload) {
	m.Called(payload)
}

func (m *MockValidateUser) ValidateUpdatePayload(payload *entity.UpdateUserPayload) {
	m.Called(payload)
}

// MockToken is a mock implementation of the Token interface for testing purposes.
type MockToken struct {
	mock.Mock
}

func (m *MockToken) CreateToken(userToken *entity.User, ttl time.Duration, privateKey string) *entity.TokenDetail {
	args := m.Called(userToken, ttl, privateKey)
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

type MockFileUpload struct {
	mock.Mock
}

func (m *MockFileUpload) UploadFile(buffer []byte, extension string) string {
	args := m.Called(buffer, extension)

	return args.String(0)
}

func (m *MockFileUpload) GetFile(fileName string) []byte {
	args := m.Called(fileName)

	return args.Get(0).([]byte)
}

func (m *MockFileUpload) RemoveFile(oldFileLink string) {
	m.Called(oldFileLink)
}

type MockFileProcessing struct {
	mock.Mock
}

func (m *MockFileProcessing) CompressImage(buffer []byte, to file_handling.ConvertTo) ([]byte, string) {
	args := m.Called(buffer, to)

	return args.Get(0).([]byte), args.String(1)
}

func (m *MockFileProcessing) AddWatermark(buffer []byte) []byte {
	args := m.Called(buffer)

	return args.Get(0).([]byte)
}
