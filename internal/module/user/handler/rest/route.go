package rest

import "github.com/gofiber/fiber/v2"

func (h *userHandler) UserRoute(router fiber.Router) {
	router.Post("/register", h.register)
	router.Post("/otp", h.verification)
	router.Post("/new-otp", h.sendOtpNumberVerification)
	router.Post("/login", h.login)
	router.Delete("/logout", h.middleware.UserBearer, h.logout)
	router.Get("/current", h.middleware.UserBearer, h.getCurrentUser)
	router.Patch("/refresh-token", h.middleware.UserRefreshBearer, h.refreshToken)
	router.Post("/forgot-password", h.forgotPassword)
	router.Post("/reset-password", h.resetPassword)
}

func (h *superAdminHandler) SuperAdminRoute(router fiber.Router) {
	router.Post("/roles/create", h.middleware.UserBearer, h.middleware.RBACMiddleware("create", "roles"), h.createRolePermission)
}
