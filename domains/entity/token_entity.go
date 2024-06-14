package entity

// TokenDetail struct represents the details of a token.
// This struct is used to encapsulate information about a generated token, including its value, associated user, and expiration details.
type TokenDetail struct {
	// Token Info
	Token     string // Token the actual access token or refresh token
	TokenId   string // TokenId used for cache validation
	ExpiresIn int64  // ExpiresIn the duration in seconds until the token expires
	MaxAge    int    // MaxAge the maximum age of the token in seconds

	// User Info, Can add more such as username, role, etc.
	UserId string // UserId the ID of the user to whom the token belongs
}
