package repository_test

import (
	"github.com/wisle25/be-template/domains/entity"
	"github.com/wisle25/be-template/infrastructures/services"
	"github.com/wisle25/be-template/tests/db_helper"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/infrastructures/generator"
	"github.com/wisle25/be-template/infrastructures/repository"
)

func TestUserRepository(t *testing.T) {
	// Arrange
	config := commons.LoadConfig("../..")
	db := services.ConnectDB(config)
	userHelperDb := &db_helper.UserHelperDB{
		DB: db,
	}
	defer userHelperDb.CleanUserDB()

	uuidGenerator := generator.NewUUIDGenerator()
	userRepositoryPG := repository.NewUserRepositoryPG(db, uuidGenerator)

	// Arrange
	payload := &entity.RegisterUserPayload{
		Username: "uname",
		Password: "password",
		Email:    "hand@gmail.com",
	}

	t.Run("AddUser", func(t *testing.T) {
		// Action
		responseId := userRepositoryPG.AddUser(payload)

		// Assert
		usersList := userHelperDb.GetUsers()

		assert.NotNil(t, responseId)
		assert.Equal(t, payload.Username, usersList[0].Username)
		assert.Equal(t, payload.Email, usersList[0].Email)
		assert.Equal(t, len(usersList), 1)
	})

	t.Run("GetUserForLogin", func(t *testing.T) {
		t.Run("Should get user by email", func(t *testing.T) {
			// Action and Assert
			assert.NotPanics(t, func() {
				userRepositoryPG.GetUserForLogin("hand@gmail.com")
			})
		})

		t.Run("Should get user by username", func(t *testing.T) {
			// Action and Assert
			assert.NotPanics(t, func() {
				userRepositoryPG.GetUserForLogin("uname")
			})
		})

		t.Run("Should raise panic if user is not existed", func(t *testing.T) {
			// Action and Assert
			assert.PanicsWithError(t, "User not found!", func() {
				userRepositoryPG.GetUserForLogin("nonExistedUname")
			})
		})
	})

	t.Run("GetUserById", func(t *testing.T) {
		expectedUser := &userHelperDb.GetUsers()[0]

		t.Run("Should get user by id", func(t *testing.T) {
			// Arrange
			var user *entity.User

			// Actions
			assert.NotPanics(t, func() {
				user = userRepositoryPG.GetUserById(expectedUser.Id)
			})
			assert.Equal(t, expectedUser, user)
		})

		t.Run("Should raise panic if user is not existed", func(t *testing.T) {
			assert.PanicsWithError(t, "User not found!", func() {
				userRepositoryPG.GetUserById(uuidGenerator.Generate())
			})
		})
	})
}
