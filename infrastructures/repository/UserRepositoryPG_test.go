package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/domains/users"
	"github.com/wisle25/be-template/infrastructures/database"
	"github.com/wisle25/be-template/infrastructures/database/db_helper"
	"github.com/wisle25/be-template/infrastructures/generator"
	"github.com/wisle25/be-template/infrastructures/repository"
)

func TestUserRepository(t *testing.T) {
	config := commons.LoadConfig("../../")
	db := database.ConnectDB(config)
	userHelperDb := &db_helper.UserHelperDB{
		DB: db,
	}
	defer userHelperDb.DB.Close()

	t.Run("VerifyUsername", func(t *testing.T) {
		uuidGenerator := generator.NewUUIDGenerator()
		userRepositoryPG := repository.NewUserRepositoryPG(db, uuidGenerator)

		t.Run("Should raising panic when username is not available", func(t *testing.T) {
			defer userHelperDb.CleanUserDB()

			// Arrange
			userHelperDb.AddUserDB(&users.RegisterUserPayload{
				Username: "uname",
				Password: "password",
				Email:    "hand@gmail.com",
			})

			// Action and Assert
			assert.PanicsWithError(t, "Username is already in use!", func() {
				userRepositoryPG.VerifyUsername("uname")
			})
		})

		t.Run("Not raising panic when username is available", func(t *testing.T) {
			// Action and Assert
			assert.NotPanics(t, func() {
				userRepositoryPG.VerifyUsername("uname")
			})
		})
	})

	t.Run("AddUser", func(t *testing.T) {
		defer userHelperDb.CleanUserDB()

		uuidGenerator := generator.NewUUIDGenerator()
		userRepositoryPG := repository.NewUserRepositoryPG(db, uuidGenerator)

		// Arrange
		payload := &users.RegisterUserPayload{
			Username: "uname",
			Password: "password",
			Email:    "hand@gmail.com",
		}

		// Action
		responseId := userRepositoryPG.AddUser(payload)
		usersList := userHelperDb.GetUsers()

		// Assert
		assert.NotNil(t, responseId, "Id shouldn't be nil!")
		assert.Equal(t, payload.Username, usersList[0].Username, "Username should be equal!")
		assert.Equal(t, payload.Email, usersList[0].Email, "Email should be equal!")
		assert.Equal(t, len(usersList), 1)
	})

	t.Run("GetUserByIdentity", func(t *testing.T) {
		uuidGenerator := generator.NewUUIDGenerator()
		userRepositoryPG := repository.NewUserRepositoryPG(db, uuidGenerator)

		payload := &users.RegisterUserPayload{
			Username: "uname",
			Password: "password",
			Email:    "hand@gmail.com",
		}
		userHelperDb.AddUserDB(payload)
		defer userHelperDb.CleanUserDB()

		t.Run("Should get user by email", func(t *testing.T) {
			// Arrange

			// Action
			response, password := userRepositoryPG.GetUserByIdentity("hand@gmail.com")

			// Assert
			assert.NotNil(t, response)
			assert.NotNil(t, password)
			assert.Equal(t, response.Username, payload.Username)
			assert.Equal(t, response.Email, payload.Email)
		})

		t.Run("Should get user by username", func(t *testing.T) {
			// Action
			response, password := userRepositoryPG.GetUserByIdentity("uname")

			// Assert
			assert.NotNil(t, response)
			assert.NotNil(t, password)
			assert.Equal(t, response.Username, payload.Username)
			assert.Equal(t, response.Email, payload.Email)
		})

		t.Run("Should raise panic if user is not existed", func(t *testing.T) {
			// Action and Assert
			assert.PanicsWithError(t, "User not found!", func() {
				userRepositoryPG.GetUserByIdentity("nonExistedUname")
			})
		})
	})
}
