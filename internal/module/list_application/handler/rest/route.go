package rest

import "github.com/gofiber/fiber/v2"

func (h *applicationHandler) ApplicationRoute(userRouter, superadminRouter fiber.Router) {
	superadminRouter.Get("/application", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "permissions"), h.getListApplication)
	superadminRouter.Get("/permissions", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "permissions"), h.getListPermissionByApp)
}
