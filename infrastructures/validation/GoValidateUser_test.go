package validation_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/wisle25/be-template/domains/users"
	"github.com/wisle25/be-template/infrastructures/validation"
	"testing"
)

func TestUsernameValidation(t *testing.T) {
	goValidator, trans := validation.NewValidator()
	validateUser := validation.NewValidateUser(goValidator, trans)

	t.Run("Should raise panic when reaches the max limit", func(t *testing.T) {
		// Arrange
		payload := users.RegisterUserPayload{
			Username: "awekrhnawklejrklewhjrtleksjrtklerwntklrsjgklrdjgtyklrdjgtyklrdjgkltydejrt",
			Password: "Helkdaoskd",
			Email:    "aiwerkajwe@gmail.com",
		}

		// Action and Assert
		assert.Panics(t, func() {
			validateUser.ValidatePayload(payload)
		})
	})

	t.Run("Should raise panic when username is empty", func(t *testing.T) {
		payload := users.RegisterUserPayload{
			Username: "",
			Password: "Helkdaoskd",
			Email:    "aiwerkajwe@gmail.com",
		}

		// Action and Assert
		assert.Panics(t, func() {
			validateUser.ValidatePayload(payload)
		})
	})

	t.Run("Should raise panic when username contains illegal chars", func(t *testing.T) {
		payload := users.RegisterUserPayload{
			Username: "AS awdawd",
			Password: "Helkdaoskd",
			Email:    "aiwerkajwe@gmail.com",
		}

		// Action and Assert
		assert.Panics(t, func() {
			validateUser.ValidatePayload(payload)
		})
	})

	t.Run("Shouldn't be no error when everyhing is met", func(t *testing.T) {
		payload := users.RegisterUserPayload{
			Username: "myusername",
			Password: "mypassword",
			Email:    "aiwerkajwe@gmail.com",
		}

		// Action and Assert
		assert.NotPanics(t, func() {
			validateUser.ValidatePayload(payload)
		}, "Everything is met on username!")
	})
}

func TestPasswordValidation(t *testing.T) {
	goValidator, trans := validation.NewValidator()
	validateUser := validation.NewValidateUser(goValidator, trans)

	t.Run("Should raise panic when password is too short", func(t *testing.T) {
		payload := users.RegisterUserPayload{
			Username: "myusername",
			Password: "short",
			Email:    "aiwerkajwe@gmail.com",
		}
		assert.Panics(t, func() {
			validateUser.ValidatePayload(payload)
		})
	})

	t.Run("Should raise panic when password is empty", func(t *testing.T) {
		payload := users.RegisterUserPayload{
			Username: "myusername",
			Password: "",
			Email:    "aiwerkajwe@gmail.com",
		}
		assert.Panics(t, func() {
			validateUser.ValidatePayload(payload)
		})
	})

	t.Run("Shouldn't raise error when password is valid", func(t *testing.T) {
		payload := users.RegisterUserPayload{
			Username: "myusername",
			Password: "validpassword",
			Email:    "aiwerkajwe@gmail.com",
		}
		assert.NotPanics(t, func() {
			validateUser.ValidatePayload(payload)
		}, "Password is valid!")
	})
}

func TestEmailValidation(t *testing.T) {
	goValidator, trans := validation.NewValidator()
	validateUser := validation.NewValidateUser(goValidator, trans)

	t.Run("Should raise panic when email is invalid", func(t *testing.T) {
		payload := users.RegisterUserPayload{
			Username: "myusername",
			Password: "validpassword",
			Email:    "invalid-email",
		}
		assert.Panics(t, func() {
			validateUser.ValidatePayload(payload)
		})
	})

	t.Run("Should raise panic when email is empty", func(t *testing.T) {
		payload := users.RegisterUserPayload{
			Username: "myusername",
			Password: "validpassword",
			Email:    "",
		}
		assert.Panics(t, func() {
			validateUser.ValidatePayload(payload)
		})
	})

	t.Run("Shouldn't raise error when email is valid", func(t *testing.T) {
		payload := users.RegisterUserPayload{
			Username: "myusername",
			Password: "validpassword",
			Email:    "valid.email@example.com",
		}
		assert.NotPanics(t, func() {
			validateUser.ValidatePayload(payload)
		}, "Email is valid!")
	})
}
