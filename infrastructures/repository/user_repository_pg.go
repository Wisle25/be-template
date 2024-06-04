package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/be-template/applications/generator"
	"github.com/wisle25/be-template/domains/entity"
	"github.com/wisle25/be-template/domains/repository"
)

type UserRepositoryPG struct {
	db          *sql.DB
	idGenerator generator.IdGenerator
}

func NewUserRepositoryPG(db *sql.DB, idGenerator generator.IdGenerator) repository.UserRepository {
	return &UserRepositoryPG{
		db:          db,
		idGenerator: idGenerator,
	}
}

func (r *UserRepositoryPG) AddUser(payload *entity.RegisterUserPayload) string {
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
		panic(fmt.Errorf("user_repo_pg_error: add user: %v", err))
	}

	return returnedId
}

func (r *UserRepositoryPG) VerifyUsername(username string) {
	var id string

	// Query
	query := "SELECT id FROM users WHERE username = $1"
	err := r.db.QueryRow(query, username).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		} else {
			panic(fmt.Errorf("user_repo_pg_error: verify username: %v", err))
		}
	}

	panic(fiber.NewError(fiber.StatusConflict, "Username is already in use!"))
}

func (r *UserRepositoryPG) GetUserByIdentity(identity string) (string, string) {
	var userId string
	var encryptedPassword string

	// Query
	query := "SELECT id, password FROM users WHERE email = $1 OR username = $1"
	err := r.db.QueryRow(query, identity).Scan(&userId, &encryptedPassword)

	// Evaluate
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			panic(fiber.NewError(fiber.StatusNotFound, "User not found!"))
		} else {
			panic(fmt.Errorf("user_repo_pg_error: get userId by identity %v", err))
		}
	}

	return userId, encryptedPassword
}
