package repository

import (
	"database/sql"
	"fmt"
	"github.com/wisle25/be-template/applications/generator"
	"github.com/wisle25/be-template/domains/users"
	"github.com/wisle25/be-template/infrastructures/database"
)

type UserRepositoryPG struct {
	db          *sql.DB
	idGenerator generator.IdGenerator
}

func NewUserRepositoryPG(db *sql.DB, idGenerator generator.IdGenerator) *UserRepositoryPG {
	return &UserRepositoryPG{
		db:          db,
		idGenerator: idGenerator,
	}
}

func (r *UserRepositoryPG) AddUser(payload *users.RegisterUserPayload) string {
	// Create ID
	id := r.idGenerator.Generate()

	// Query
	query := `INSERT INTO 
    			users(id, username, password, email) 
			  VALUES
			      ($1, $2, $3, $4)
			  RETURNING id`

	var returnedId string
	err := r.db.QueryRow(
		query,
		id,
		payload.Username,
		payload.Password,
		payload.Email,
	).Scan(&returnedId)

	if err != nil {
		panic(err)
	}

	return returnedId
}

func (r *UserRepositoryPG) VerifyUsername(username string) {
	// Query
	query := "SELECT id, username, email FROM users WHERE username = $1"
	result, err := r.db.Query(query, username)

	if err != nil {
		panic(err)
	}

	rows := database.GetTableDB[users.User](result)
	fmt.Printf("Rows len: %d\n", len(rows))

	if len(rows) > 0 {
		panic(fmt.Errorf("username is already in use"))
	}
}
