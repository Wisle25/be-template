package db_helper

import (
	"github.com/wisle25/be-template/domains/users"
	"github.com/wisle25/be-template/infrastructures/database"
)

func AddUserDB(payload *users.RegisterUserPayload) {
	query := "INSERT INTO users(id, username, password, email) VALUES ($1, $2, $3, $4)"
	_, _ = database.DB.Exec(
		query,
		"user-123",
		payload.Username,
		payload.Password,
		payload.Email,
	)
}

func GetUsers() []users.User {
	query := "SELECT id FROM users"
	row := database.DB.QueryRow(query)

	return database.GetTableDB[users.User](row)
}

func CleanUserDB() {
	_, _ = database.DB.Exec("TRUNCATE TABLE users")
}
