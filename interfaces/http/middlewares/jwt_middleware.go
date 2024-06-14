package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/be-template/applications/use_case"
)

// JwtMiddleware gonna verifying user credential for authentication/authorization reason
// Using UserUseCase since it handles both authentication and users itself
type JwtMiddleware struct {
	userUseCase *use_case.UserUseCase
}

func NewJwtMiddleware(userUseCase *use_case.UserUseCase) *JwtMiddleware {
	return &JwtMiddleware{userUseCase}
}

func (m *JwtMiddleware) GuardJWT(c *fiber.Ctx) error {
	// Getting access token
	accessToken := c.Cookies("access_token")

	if accessToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "You are not logged in!")
	}

	userId, accessTokenDetail := m.userUseCase.ExecuteGuard(accessToken)

	if userId == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Session invalid or expired!")
	}

	// Add additional information
	c.Locals("access_token_id", accessTokenDetail.TokenId)
	c.Locals("user_id", userId.(string))

	return c.Next()
}
