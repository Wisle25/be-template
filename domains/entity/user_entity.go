package entity

type RegisterUserPayload struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	Email           string `json:"email"`
	ConfirmPassword string `json:"confirmPassword"`
}

type LoginUserPayload struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type User struct {
	Id       string
	Username string
	Email    string
}
