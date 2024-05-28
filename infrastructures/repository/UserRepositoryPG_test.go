package repository_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/domains/users"
	"github.com/wisle25/be-template/infrastructures/database"
	"github.com/wisle25/be-template/infrastructures/database/db_helper"
	"github.com/wisle25/be-template/infrastructures/generator"
	"github.com/wisle25/be-template/infrastructures/repository"
	"testing"
)

func PrepareRepositoryTest() {
	commons.LoadConfig("../../.")
	database.ConnectDB(commons.Cfg)
}

func DisconnectDB() {
	_ = database.DB.Close()
}

func TestVerifyUsername(t *testing.T) {
	PrepareRepositoryTest()
	defer DisconnectDB()

	uuidGenerator := &generator.UUIDGenerator{}

	t.Run("raise panic when username is not available", func(t *testing.T) {
		defer db_helper.CleanUserDB()

		// Arrange
		db_helper.AddUserDB(&users.RegisterUserPayload{
			Username: "uname",
			Password: "password",
			Email:    "hand@gmail.com",
		})

		userRepositoryPG := repository.NewUserRepositoryPG(database.DB, uuidGenerator)

		// Assert
		assert.PanicsWithError(t, "Username is not available!", func() {
			userRepositoryPG.VerifyUsername("uname")
		})
	})

	t.Run("Not raising panic when username is available", func(t *testing.T) {
		defer db_helper.CleanUserDB()

		// Arrange
		userRepositoryPG := repository.NewUserRepositoryPG(database.DB, uuidGenerator)

		// Assert
		assert.NotPanics(t, func() {
			userRepositoryPG.VerifyUsername("uname")
		})
	})
}

func TestAddUser(t *testing.T) {
	PrepareRepositoryTest()
	defer DisconnectDB()

	uuidGenerator := &generator.UUIDGenerator{}

	t.Run("Should really add user to database", func(t *testing.T) {
		defer db_helper.CleanUserDB()

		// Arrange
		payload := &users.RegisterUserPayload{
			Username: "uname",
			Password: "password",
			Email:    "hand@gmail.com",
		}
		userRepositoryPG := repository.NewUserRepositoryPG(database.DB, uuidGenerator)

		// Action
		response := userRepositoryPG.AddUser(payload)
		usersList := db_helper.GetUsers()

		// Assert
		assert.NotNil(t, response.Id, "Id shouldn't be nil!")
		assert.Equal(t, payload.Username, response.Username, "Username should be equal!")
		assert.Equal(t, payload.Email, response.Email, "Email should be equal!")
		assert.Equal(t, len(usersList), 1)
	})
}
