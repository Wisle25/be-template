﻿package security

import (
	"github.com/wisle25/be-template/domains/tokens"
	"time"
)

// Token interface defines methods for creating and validating tokens.
type Token interface {
	// CreateToken generates a new token for the given user ID with a specified time-to-live (ttl) duration.
	// The token is signed using the provided private key.
	// Returns a TokenDetail which contains the token and its metadata.
	CreateToken(userID string, ttl time.Duration, privateKey string) *tokens.TokenDetail

	// ValidateToken validates the given token using the provided public key.
	// Returns a TokenDetail which contains the token's metadata if the token is valid.
	ValidateToken(token string, publicKey string) *tokens.TokenDetail
}
