package repository_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/wisle25/be-template/domains/users"
	"github.com/wisle25/be-template/infrastructures/database/db_test"
	"github.com/wisle25/be-template/infrastructures/repository"
	"testing"
)

func TestVerifyUsername(t *testing.T) {
	t.Run("raise panic when username is not available", func(t *testing.T) {
		defer db_test.CleanUserDB()

		// Arrange
		db_test.AddUserDB(&users.RegisterUserPayload{
			Username: "uname",
			Password: "password",
			Email:    "hand@gmail.com",
		})

		userRepositoryPG := repository.NewUserRepositoryPG()

		// Assert
		assert.PanicsWithError(t, "Username is not available!", func() {
			userRepositoryPG.VerifyUsername("uname")
		})
	})

	t.Run("Not raising panic when username is available", func(t *testing.T) {
		defer db_test.CleanUserDB()

		// Arrange
		userRepositoryPG := repository.NewUserRepositoryPG()

		// Assert
		assert.NotPanics(t, func() {
			userRepositoryPG.VerifyUsername("uname")
		})
	})
}

func TestAddUser(t *testing.T) {
	t.Run("Should really add user to database", func(t *testing.T) {
		defer db_test.CleanUserDB()

		// Arrange
		payload := &users.RegisterUserPayload{
			Username: "uname",
			Password: "password",
			Email:    "hand@gmail.com",
		}
		userRepositoryPG := repository.NewUserRepositoryPG()

		// Action
		response := userRepositoryPG.AddUser(payload)

		assert.Equal(t, users.RegisterUserResponse{})
	})
}
