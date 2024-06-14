package repository

import "github.com/wisle25/be-template/domains/entity"

// UserRepository defines methods for interacting with the user-related data in the database.
type UserRepository interface {
	// AddUser adds a new user to the database using the provided registration payload.
	// It should raise panic if email or username is already taken.
	// Returns the ID of the newly created user.
	AddUser(payload *entity.RegisterUserPayload) string

	// GetUserForLogin retrieves user details based on the provided identity (username or email) for login purpose.
	// It should raise panic if user is not existed
	// Returns the user's ID and hashed password to be decrypted.
	GetUserForLogin(identity string) (string, string)

	// GetUserById Get detailed information about User
	// It should raise panic if user is not existed
	GetUserById(id string) *entity.User
}
