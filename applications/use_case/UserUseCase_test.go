package use_case_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wisle25/be-template/applications/use_case"
	"github.com/wisle25/be-template/domains/users"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) VerifyUsername(username string) {
	args := m.Called(username)

	if args.Get(0) != nil {
		panic(args.Get(0).(error))
	}
}

func (m *MockUserRepository) AddUser(user *users.RegisterUserPayload) *users.RegisterUserResponse {
	args := m.Called(user)

	return args.Get(0).(*users.RegisterUserResponse)
}

type MockPasswordHash struct {
	mock.Mock
}

func (m *MockPasswordHash) Hash(password string) string {
	args := m.Called(password)

	return args.String(0)
}

type MockValidator struct {
    mock.Mock
}

func (m* MockValidator) ValidatePayload(s interface{}) {
    args := m.Called(s)

    if args.Get(0) != nil {
        panic(args.Get(0).(error))
    }
}

func TestAddUserUseCase(t *testing.T) {
	t.Run("should orchestrate the add user action correctly", func(t *testing.T) {
		// Arrange
		useCasePayload := &users.RegisterUserPayload{
			Username: "user",
			Password: "secret",
			Email:    "user@example.com",
		}
		mockRegisteredUser := &users.RegisterUserResponse{
			Id:       "user-123",
			Username: useCasePayload.Username,
			Email:    useCasePayload.Email,
		}

		mockUserRepository := new(MockUserRepository)
		mockPasswordHash := new(MockPasswordHash)
        mockValidator := new(MockValidator)

		mockUserRepository.On("VerifyUsername", useCasePayload.Username).Return(nil)
		mockPasswordHash.On("Hash", useCasePayload.Password).Return("encrypted_password")
		mockUserRepository.On("AddUser", mock.MatchedBy(func(payload *users.RegisterUserPayload) bool {
			return payload.Username == useCasePayload.Username && payload.Email == useCasePayload.Email && payload.Password == "encrypted_password"
		})).Return(mockRegisteredUser)
        mockValidator.On("ValidatePayload", useCasePayload).Return(nil)
        
		// Action
		addUserUseCase := use_case.NewAddUserUseCase(mockUserRepository, mockPasswordHash, mockValidator)
		registeredUser := addUserUseCase.Execute(useCasePayload)

		// Assert
		assert.Equal(t, mockRegisteredUser, registeredUser)
        mockValidator.AssertCalled(t, "ValidatePayload", useCasePayload)
		mockUserRepository.AssertCalled(t, "VerifyUsername", useCasePayload.Username)
		mockPasswordHash.AssertCalled(t, "Hash", "secret")
		mockUserRepository.AssertCalled(t, "AddUser", mock.MatchedBy(func(payload *users.RegisterUserPayload) bool {
			return payload.Username == useCasePayload.Username && payload.Email == useCasePayload.Email && payload.Password == "encrypted_password"
		}))
	})

	t.Run("should panic when username is not available", func(t *testing.T) {
		// Arrange
		useCasePayload := &users.RegisterUserPayload{
			Username: "user",
			Password: "secret",
			Email:    "user@example.com",
		}

		mockUserRepository := new(MockUserRepository)
		mockPasswordHash := new(MockPasswordHash)
        mockValidator := new(MockValidator)

		expectedPanic := "USERNAME_NOT_AVAILABLE"
		mockUserRepository.On("VerifyUsername", useCasePayload.Username).Return(fmt.Errorf(expectedPanic))
        mockValidator.On("ValidatePayload", useCasePayload).Return(nil)

		addUserUseCase := use_case.NewAddUserUseCase(mockUserRepository, mockPasswordHash, mockValidator)

		// Action and Assert
		assert.PanicsWithError(t, expectedPanic, func() {
			addUserUseCase.Execute(useCasePayload)
		})
		mockUserRepository.AssertCalled(t, "VerifyUsername", useCasePayload.Username)
	})

	t.Run("should panic when validation fail", func(t *testing.T) {
		// Arrange
		useCasePayload := &users.RegisterUserPayload{
			Username: "user",
			Password: "secret",
		}

		mockUserRepository := new(MockUserRepository)
		mockPasswordHash := new(MockPasswordHash)
        mockValidator := new(MockValidator)

		expectedPanic := "VALIDATION_FAIL"
		mockUserRepository.On("VerifyUsername", useCasePayload.Username).Return(nil)
        mockValidator.On("ValidatePayload", useCasePayload).Return(fmt.Errorf(expectedPanic))

		addUserUseCase := use_case.NewAddUserUseCase(mockUserRepository, mockPasswordHash, mockValidator)

		// Action and Assert
		assert.PanicsWithError(t, expectedPanic, func() {
			addUserUseCase.Execute(useCasePayload)
		})
		mockValidator.AssertCalled(t, "ValidatePayload", useCasePayload)
	})
}
