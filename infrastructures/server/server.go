﻿package server

import (
	"errors"
	middlewares2 "github.com/wisle25/be-template/interfaces/middlewares"
	"github.com/wisle25/be-template/interfaces/users"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/infrastructures/cache"
	"github.com/wisle25/be-template/infrastructures/container"
	"github.com/wisle25/be-template/infrastructures/file_handling"
	"github.com/wisle25/be-template/infrastructures/generator"
	"github.com/wisle25/be-template/infrastructures/services"
)

func errorHandling(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	status := "error"
	code := fiber.StatusInternalServerError
	message := err.Error()

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		status = "fail"
		code = e.Code
		message = e.Message
	}

	// Send custom error
	return c.Status(code).JSON(fiber.Map{
		"status":  status,
		"message": message,
	})
}

func CreateServer(config *commons.Config) *fiber.App {
	// Load Services
	db := services.ConnectDB(config)
	redis := services.ConnectRedis(config)
	minio, bucketName := services.NewMinio(config)

	// Server
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandling,
	})

	// Middlewares
	app.Use(middlewares2.NewPlatformMiddleware)
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.ClientOrigin,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "POST,GET,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
	}))

	// Global Dependencies
	redisCache := cache.NewRedisCache(redis)
	uuidGenerator := generator.NewUUIDGenerator()
	validation := services.NewValidation()
	minioFileUpload := file_handling.NewMinioFileUpload(minio, uuidGenerator, bucketName)
	vipsFileProcessing := file_handling.NewVipsFileProcessing()

	// Use Cases
	userUseCase := container.NewUserContainer(
		config,
		db,
		redisCache,
		uuidGenerator,
		vipsFileProcessing,
		minioFileUpload,
		validation,
	)

	// Custom Middleware
	jwtMiddleware := middlewares2.NewJwtMiddleware(userUseCase)

	// Router
	users.NewUserRouter(app, jwtMiddleware, userUseCase)

	return app
}
