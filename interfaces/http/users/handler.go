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
