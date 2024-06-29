package middlewares

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/be-template/applications/use_case"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/domains/entity"
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
	var accessToken string

	if c.Locals("isMobile").(bool) {
		accessToken = strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	} else {
		accessToken = c.Cookies("access_token") // From cookies

	}

	if accessToken == "" {
		commons.ThrowClientError(fiber.StatusUnauthorized, "You are not logged in!")
	}

	// Verifying access token
	userInfoJSON, accessTokenDetail := m.userUseCase.ExecuteGuard(accessToken)

	if userInfoJSON == nil {
		commons.ThrowClientError(fiber.StatusUnauthorized, "Session invalid or expired!")
	}

	// Unmarshal userInfo JSON
	var userInfo entity.User
	err := json.Unmarshal([]byte(userInfoJSON.(string)), &userInfo)
	if err != nil {
		commons.ThrowServerError("refresh_token_err: unable to unmarshal json user info: %v", err)
	}

	// Add additional information
	c.Locals("accessTokenId", accessTokenDetail.TokenId)
	c.Locals("userInfo", userInfo)

	return c.Next()
}
