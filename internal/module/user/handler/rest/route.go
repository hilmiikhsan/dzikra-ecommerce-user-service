package rest

import "github.com/gofiber/fiber/v2"

func (h *userHandler) UserRoute(router fiber.Router) {
	router.Post("/register", h.register)
	router.Post("/otp", h.verification)
	router.Post("/new-otp", h.sendOtpNumberVerification)
	router.Post("/login", h.login)
	router.Delete("/logout", h.middleware.UserBearer, h.logout)
}
