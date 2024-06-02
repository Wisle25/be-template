package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/be-template/applications/use_case"
)

func NewUserRouter(app *fiber.App, useCase *use_case.UserUseCase) {
	userHandler := NewUserHandler(useCase)

	app.Post("/users", userHandler.AddUser)
	app.Post("/auths", userHandler.Login)
	app.Put("/auths", userHandler.RefreshToken)
}
