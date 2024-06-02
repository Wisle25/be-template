package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/be-template/applications/generator"
	"github.com/wisle25/be-template/domains/users"
	"github.com/wisle25/be-template/infrastructures/database"
)

type UserRepositoryPG struct {
	db          *sql.DB
	idGenerator generator.IdGenerator
}

func NewUserRepositoryPG(db *sql.DB, idGenerator generator.IdGenerator) users.UserRepository {
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
		panic(fmt.Errorf("user_repo_pg_error: %v", err))
	}

	return returnedId
}

func (r *UserRepositoryPG) VerifyUsername(username string) {
	// Query
	query := "SELECT id, username, email FROM users WHERE username = $1"
	result, err := r.db.Query(query, username)

	if err != nil {
		panic(fmt.Errorf("user_repo_pg_error: %v", err))
	}

	rows := database.GetTableDB[users.User](result)

	if len(rows) > 0 {
		panic(fiber.NewError(fiber.StatusConflict, "Username is already in use!"))
	}
}

func (r *UserRepositoryPG) GetUserByIdentity(identity string) (*users.User, string) {
	var user users.User
	var encryptedPassword string

	// Query
	query := "SELECT id, username, email, password FROM users WHERE email = $1 OR username = $1"
	err := r.db.QueryRow(query, identity).Scan(&user.Id, &user.Username, &user.Email, &encryptedPassword)

	// Evaluate
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			panic(fiber.NewError(fiber.StatusNotFound, "User not found!"))
		} else {
			panic(fmt.Errorf("user_repo_pg_error: %v", err))
		}
	}

	return &user, encryptedPassword
}
