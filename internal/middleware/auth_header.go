package middleware

import (
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (m *UserMiddleware) UserBearer(c *fiber.Ctx) error {
	accessToken := c.Get(constants.HeaderAuthorization)

	// If the cookie is not set, return an unauthorized status
	if accessToken == "" {
		log.Error().Msg("middleware::UserBearer - Unauthorized [Header not set]")
		c.Status(fiber.StatusUnauthorized).JSON(response.Error(constants.ErrAccessTokenIsRequired))
	}

	// remove the Bearer prefix
	if len(accessToken) > 7 {
		accessToken = accessToken[7:]
	}

	// Parse the JWT string and store the result in `claims`
	claims, err := m.jwt.ParseMiddlewareTokenString(c.Context(), accessToken)
	if err != nil {
		log.Error().Err(err).Any("payload", accessToken).Msg("middleware::UserBearer - Error while parsing token")
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(constants.ErrTokenAlreadyExpired))
	}

	c.Locals("user_id", claims.UserID)
	c.Locals("username", claims.Username)
	c.Locals("email", claims.Email)
	c.Locals("full_name", claims.FullName)

	// If the token is valid, pass the request to the next handler
	return c.Next()
}
