package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/be-template/applications/use_case"
	"github.com/wisle25/be-template/domains/users"
)

type UserHandler struct {
	useCase *use_case.UserUseCase
}

func NewUserHandler(useCase *use_case.UserUseCase) *UserHandler {
	return &UserHandler{
		useCase: useCase,
	}
}

func (h *UserHandler) AddUser(c *fiber.Ctx) error {
	var payload users.RegisterUserPayload
	_ = c.BodyParser(&payload)
	returnedId := h.useCase.ExecuteAdd(&payload)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   returnedId,
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var payload users.LoginUserPayload
	_ = c.BodyParser(&payload)
	accessTokenDetail, refreshTokenDetail := h.useCase.ExecuteLogin(&payload)

	// Insert the tokens to cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessTokenDetail.Token,
		Path:     "/",
		MaxAge:   accessTokenDetail.MaxAge,
		Secure:   true,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenDetail.Token,
		Path:     "/",
		MaxAge:   accessTokenDetail.MaxAge,
		Secure:   true,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   accessTokenDetail.MaxAge,
		Secure:   false,
		HTTPOnly: false,
		Domain:   "localhost",
	})

	// ...
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully logged in!",
	})
}
