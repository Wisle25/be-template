package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/be-template/infrastructures/container"
)

func NewUserRouter(app *fiber.App) {
	userHandler := NewUserHandler(container.NewUserContainer())

	app.Post("/users", userHandler.AddUser)
}
