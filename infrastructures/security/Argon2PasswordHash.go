package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/wisle25/be-template/applications/security"

	"golang.org/x/crypto/argon2"
)

// Argon2PasswordHash Using argon2 for password hashing
type Argon2PasswordHash struct /* implements PasswordHash */ {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

func NewArgon2() security.PasswordHash {
	return &Argon2PasswordHash{
		Time:    3,
		Memory:  64 * 1024,
		Threads: 2,
		SaltLen: 16,
		KeyLen:  32,
	}
}

// Hash the password
func (a *Argon2PasswordHash) Hash(password string) string {
	salt, err := a.GenerateSalt()
	if err != nil {
		panic(err)
	}

	hash := argon2.IDKey([]byte(password), salt, a.Time, a.Memory, a.Threads, a.KeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	result := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, a.Memory, a.Time, a.Threads, b64Salt, b64Hash)

	return result
}

func (a *Argon2PasswordHash) GenerateSalt() ([]byte, error) {
	secret := make([]byte, a.SaltLen)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// func ComparePassword(password, hashedPassword string) bool {
// 	c, salt, hash, err := DecodeHash(hashedPassword)
// 	if err != nil {
// 		panic(err)
// 	}

// 	otherHash := argon2.IDKey([]byte(password), salt, c.Time, c.Memory, c.Threads, c.KeyLen)

// 	return subtle.ConstantTimeCompare(hash, otherHash) == 1
// }

// func DecodeHash(hashedPassword string) (c *Argon2Hash, salt, hash []byte, err error) {
// 	vals := strings.Split(hashedPassword, "$")
// 	if len(vals) != 6 {
// 		return nil, nil, nil, errors.New("password is not in the correct format")
// 	}

// 	var version int
// 	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
// 	if err != nil || version != argon2.Version {
// 		return nil, nil, nil, err
// 	}

// 	c = &Argon2Hash{}
// 	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &c.Memory, &c.Time, &c.Threads)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}

// 	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	c.SaltLen = uint32(len(salt))

// 	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	c.KeyLen = uint32(len(hash))

// 	return c, salt, hash, nil
// }
