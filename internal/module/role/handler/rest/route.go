package rest

import "github.com/gofiber/fiber/v2"

func (h *roleHandler) RoleRoute(userRouter, superadminRouter fiber.Router) {
	userRouter.Get("/roles/:role_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "roles"), h.getDetailRole)
	superadminRouter.Get("/roles", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "roles"), h.getListRole)
	superadminRouter.Post("/roles/create", h.middleware.UserBearer, h.middleware.RBACMiddleware("create", "roles"), h.createRolePermission)
	superadminRouter.Patch("/roles/update/:role_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("update", "roles"), h.updateRolePermission)
	superadminRouter.Delete("/roles/remove/:role_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("delete", "roles"), h.removeRolePermission)
}
