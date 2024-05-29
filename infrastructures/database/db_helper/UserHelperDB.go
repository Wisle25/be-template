package db_helper

import (
	"github.com/google/uuid"
	"github.com/wisle25/be-template/domains/users"
	"github.com/wisle25/be-template/infrastructures/database"
)

func AddUserDB(payload *users.RegisterUserPayload) {
	id, _ := uuid.NewV7()
	query := "INSERT INTO users(id, username, password, email) VALUES ($1, $2, $3, $4)"
	_, err := database.DB.Exec(
		query,
		id,
		payload.Username,
		payload.Password,
		payload.Email,
	)

	if err != nil {
		panic(err)
	}
}

func GetUsers() []users.User {
	query := "SELECT id, username, email FROM users"
	rows, _ := database.DB.Query(query)

	return database.GetTableDB[users.User](rows)
}

func CleanUserDB() {
	_, _ = database.DB.Query("TRUNCATE TABLE users")
}
