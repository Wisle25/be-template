package security

// PasswordHash interface defines methods for hashing and comparing passwords.
// This interface is typically implemented to provide password hashing functionality using algorithms such as Argon2, bcrypt, etc.
type PasswordHash interface {
	// Hash hashes the given password and returns the hashed result as a string.
	Hash(password string) string

	// Compare compares the given password with the hashed password.
	// It should panic or return an error if the comparison fails.
	Compare(password string, hashedPassword string)
}
