package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/be-template/applications/use_case"
	"github.com/wisle25/be-template/interfaces/middlewares"
)

func NewUserRouter(
	app *fiber.App,
	jwtMiddleware *middlewares.JwtMiddleware,
	useCase *use_case.UserUseCase,
) {
	userHandler := NewUserHandler(useCase)

	app.Post("/users", userHandler.RegisterUser).Name("Register")
	app.Post("/auths", userHandler.Login).Name("Login")
	app.Get("/auths", jwtMiddleware.GuardJWT, userHandler.GetLoggedUser)
	app.Put("/auths", userHandler.RefreshToken)
	app.Delete("/auths", jwtMiddleware.GuardJWT, userHandler.Logout)
	app.Get("/users", userHandler.GetAllUsers)
	app.Get("/users/:id", userHandler.GetUserById)
	app.Put("/users/:id", jwtMiddleware.GuardJWT, userHandler.UpdateUserById)
}
