package entity

import "mime/multipart"

// RegisterUserPayload represents the payload for user registration.
type RegisterUserPayload struct {
	Username        string `json:"username"`        // Username chosen by the user
	Password        string `json:"password"`        // Password chosen by the user
	Email           string `json:"email"`           // User's email address
	ConfirmPassword string `json:"confirmPassword"` // Confirmation of the user's password
}

// LoginUserPayload represents the payload for user login.
type LoginUserPayload struct {
	Identity string `json:"identity"` // User's identity which could be username or email
	Password string `json:"password"` // User's password
}

type UpdateUserPayload struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Avatar          *multipart.FileHeader
}

// User represents a user in the system.
type User struct {
	Id         string // Id for the user
	Username   string // Username of the user, Username should be unique
	Email      string // Email address of the user, Email should be unique
	AvatarLink string // AvatarLink to the user's avatar image
}
