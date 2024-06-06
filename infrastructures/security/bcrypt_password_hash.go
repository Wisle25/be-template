package security

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/be-template/applications/security"
	"golang.org/x/crypto/bcrypt"
)

// BcryptPasswordHash Using argon2 for password hashing
type BcryptPasswordHash struct /* implements PasswordHash */ {

}

func NewBcrypt() security.PasswordHash {
	return &BcryptPasswordHash{}
}

// Hash the password
func (a *BcryptPasswordHash) Hash(password string) string {
	passwordBytes := []byte(password)
	hashed, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	passwordBytes = nil

	if err != nil {
		panic(fmt.Errorf("bcrypt_err: hash: %v", err))
	}

	hashedStr := string(hashed)
	hashed = nil

	return hashedStr
}

func (a *BcryptPasswordHash) Compare(password string, hashedPassword string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			panic(fiber.NewError(fiber.StatusUnauthorized, "Password is incorrect!"))
		}

		panic(fmt.Errorf("bcrypt_err: compare: %v", err))
	}
}
