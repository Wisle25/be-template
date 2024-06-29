package commons

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func ThrowClientError(statusCode int, message string) {
	log.Printf("%s", message)
	panic(fiber.NewError(statusCode, message))
}

func ThrowServerError(message string, err error) {
	log.Printf("%s: %v", message, err)
	panic(fmt.Errorf("%s: %v", message, err))
}
