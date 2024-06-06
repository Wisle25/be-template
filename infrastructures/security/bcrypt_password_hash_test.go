package security_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wisle25/be-template/infrastructures/security"
)

func TestArgon2PasswordHash_Hash(t *testing.T) {
	t.Run("should encrypt password correctly", func(t *testing.T) {
		// Arrange
		argon2PasswordHash := security.NewBcrypt()
		plainPassword := "plain_password"

		// Act
		encryptedPassword := argon2PasswordHash.Hash(plainPassword)

		// Assert
		assert.NotEqual(t, plainPassword, encryptedPassword, "encrypted password should not be the same as plain password")
		assert.IsType(t, "string", encryptedPassword, "encrypted password should be a string")
		fmt.Printf("encryptedPassword: %s\n", encryptedPassword)
	})
}

func TestArgon2PasswordHash_Compare(t *testing.T) {
	// Arrange
	argon2PasswordHash := security.NewBcrypt()
	plainPassword := "plain_password"

	// Act
	encryptedPassword := argon2PasswordHash.Hash(plainPassword)

	t.Run("Comparison shouldnt return error when its matched", func(t *testing.T) {
		// Action and Assert
		assert.NotPanics(t, func() {
			argon2PasswordHash.Compare(plainPassword, encryptedPassword)
		})
	})

	t.Run("Comparison should raise panic when its not matched", func(t *testing.T) {
		// Action
		assert.PanicsWithError(t, "Password is incorrect!", func() {
			argon2PasswordHash.Compare("anotherPassword", encryptedPassword)
		})
	})
}
