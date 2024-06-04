package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/be-template/applications/use_case"
	"github.com/wisle25/be-template/domains/entity"
	"time"
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
	// Use Case
	var payload entity.RegisterUserPayload
	_ = c.BodyParser(&payload)
	returnedId := h.useCase.ExecuteAdd(&payload)

	// Response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   returnedId,
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	// Use Case
	var payload entity.LoginUserPayload
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

	// Response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully logged in!",
	})
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	// Use Case
	refreshToken := c.Cookies("refresh_token")

	accessTokenDetail := h.useCase.ExecuteRefreshToken(refreshToken)

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

	// Response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	// Use Case
	refreshToken := c.Cookies("refresh_token")
	accessTokenId := c.Locals("access_token_id").(string)

	h.useCase.ExecuteLogout(refreshToken, accessTokenId)
	// Remove from cookie
	expiredTime := time.Now().Add(-time.Hour * 24)

	c.Cookie(&fiber.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: expiredTime,
	})
	c.Cookie(&fiber.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Expires: expiredTime,
	})
	c.Cookie(&fiber.Cookie{
		Name:    "logged_in",
		Value:   "",
		Expires: expiredTime,
	})

	// Response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully logged out!",
	})
}
