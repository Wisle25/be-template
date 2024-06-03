// Package security provides functionalities for creating and validating JWT tokens.
package security

import (
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wisle25/be-template/applications/generator"
	"github.com/wisle25/be-template/applications/security"
	"github.com/wisle25/be-template/domains/tokens"
	"time"
)

// JwtToken struct provides methods for creating and validating token using JWT.
type JwtToken struct {
	idGenerator generator.IdGenerator
}

// NewJwtToken returns a new instance of JwtToken.
func NewJwtToken(idGenerator generator.IdGenerator) security.Token {
	return &JwtToken{
		idGenerator: idGenerator,
	}
}

// CreateToken generates a new JWT token for a given user ID and time-to-live duration.
// It uses the provided private key to sign the token.
func (jt *JwtToken) CreateToken(userID string, ttl time.Duration, privateKey string) *tokens.TokenDetail {
	now := time.Now().UTC()

	// Creating token details
	td := &tokens.TokenDetail{
		TokenID:   jt.idGenerator.Generate(),
		UserID:    userID,
		ExpiresIn: now.Add(ttl).Unix(),
		MaxAge:    int(ttl.Seconds()),
	}

	// Decode the private key
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		panic(fmt.Errorf("create_token_err: couldn't decode token private key: %v", err))
	}

	// Parse the private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		panic(fmt.Errorf("create_token_err: couldn't parse token private key: %v", err))
	}

	// Define JWT claims
	atClaims := jwt.MapClaims{
		"sub":      userID,
		"token_id": td.TokenID,
		"exp":      td.ExpiresIn,
		"iat":      now.Unix(),
		"nbf":      now.Unix(),
	}

	// Create and sign the JWT token
	td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(key)
	if err != nil {
		panic(fmt.Errorf("create_token_err: signing token error: %v", err))
	}

	return td
}

// ValidateToken verifies the given JWT token using the provided public key.
// It returns the token details if the token is valid.
func (jt *JwtToken) ValidateToken(token string, publicKey string) *tokens.TokenDetail {
	if token == "" {
		panic(fiber.NewError(
			fiber.StatusUnauthorized,
			"Session is invalid! You might be not signed in!",
		))
	}

	// Decode the public key
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		panic(fmt.Errorf("validate_token_err: couldn't decode: %v", err))
	}

	// Parse the public key
	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		panic(fmt.Errorf("validate_token_err: couldn't parse: %v", err))
	}

	// Parse and validate the JWT token
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			panic(fmt.Errorf("validate_token_err: unexpected method: %s", t.Header["alg"]))
		}

		return key, nil
	})

	if err != nil {
		panic(fmt.Errorf("validate_token_err: %v", err))
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok || !parsedToken.Valid {
		panic(fiber.NewError(fiber.StatusInternalServerError, "Invalid token!"))
	}

	// Return the token details
	return &tokens.TokenDetail{
		TokenID: fmt.Sprintf("%s", claims["token_id"]),
		UserID:  fmt.Sprintf("%s", claims["sub"]),
	}
}
