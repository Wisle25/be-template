package server

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/infrastructures/container"
	"github.com/wisle25/be-template/infrastructures/database"
	"github.com/wisle25/be-template/interfaces/http/middlewares"
	"github.com/wisle25/be-template/interfaces/http/users"
)

func errorHandling(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError
	message := err.Error()

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	// Send custom error
	return c.Status(code).JSON(fiber.Map{
		"status":  "fail",
		"message": message,
	})
}

func CreateServer(config *commons.Config) *fiber.App {
	// Load Utils
	db := database.ConnectDB(config)
	redis := database.ConnectRedis(config)

	// Server
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandling,
	})

	// Use Cases
	userUseCase := container.NewUserContainer(config, db, redis)

	// Middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "*",
		AllowCredentials: true,
	}))

	// Custom Middleware
	jwtMiddleware := middlewares.NewJwtMiddleware(userUseCase)

	// Router
	users.NewUserRouter(app, jwtMiddleware, userUseCase)

	return app
}
