package users

type RegisterUserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	Id       string
	Username string
	Email    string
}
