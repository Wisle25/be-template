package entity

// TokenDetail struct represents the details of a token.
// This struct is used to encapsulate information about a generated token, including its value, associated user, and expiration details.
type TokenDetail struct {
	Token     string
	TokenID   string
	UserID    string
	ExpiresIn int64
	MaxAge    int
}
