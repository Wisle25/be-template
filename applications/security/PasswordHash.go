package security

type PasswordHash interface {
	Hash(password string) string
}