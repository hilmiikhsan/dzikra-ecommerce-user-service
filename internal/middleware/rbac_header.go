package middleware

import (
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/adapter"
	userRole "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_role/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (m *UserMiddleware) RBACMiddleware(action, resource string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			log.Error().Msg("middleware::RBACMiddleware - user_id not found in context")
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(constants.ErrAccessTokenIsRequired))
		}

		userRolePermissions, err := userRole.NewUserRoleRepository(adapter.Adapters.DzikraPostgres).FindPermissionsByUserID(c.Context(), userID)
		if err != nil {
			log.Error().Err(err).Any("user_id", userID).Msg("middleware::RBACMiddleware - Failed to get user roles with permissions")
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error(constants.ErrInternalServerError))
		}

		requiredPermission := fmt.Sprintf("%s|%s", action, resource)
		hasPermission := false
		for _, perm := range userRolePermissions {
			if perm == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			log.Error().Any("action", action).Any("resource", resource).Msg("middleware::RBACMiddleware - User does not have permission")
			return c.Status(fiber.StatusForbidden).JSON(response.Error(constants.ErrApplicationForbidden))
		}

		return c.Next()
	}
}
