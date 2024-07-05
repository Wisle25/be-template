package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewPlatformMiddleware(c *fiber.Ctx) error {
	isMobile := false

	// Check for custom header
	if c.Get("X-Client-Type") == "mobile" {
		isMobile = true
	} else {
		// Check the User-Agent header
		userAgent := c.Get("User-Agent")
		if strings.Contains(userAgent, "Mobile") || 
		   strings.Contains(userAgent, "Android") || 
		   strings.Contains(userAgent, "iPhone") {
			isMobile = true
		}
	}

	c.Locals("isMobile", isMobile)

	return c.Next()
}