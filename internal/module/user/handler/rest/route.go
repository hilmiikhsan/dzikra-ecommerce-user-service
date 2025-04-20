package rest

import "github.com/gofiber/fiber/v2"

func (h *userHandler) UserRoute(userRouter, superadminRouter fiber.Router) {
	// auth user endpoint
	userRouter.Post("/register", h.register)
	userRouter.Post("/otp", h.verification)
	userRouter.Post("/new-otp", h.sendOtpNumberVerification)
	userRouter.Post("/login", h.login)
	userRouter.Delete("/logout", h.middleware.UserBearer, h.logout)
	userRouter.Get("/current", h.middleware.UserBearer, h.getCurrentUser)
	userRouter.Patch("/refresh-token", h.middleware.UserRefreshBearer, h.refreshToken)
	userRouter.Post("/forgot-password", h.forgotPassword)
	userRouter.Post("/reset-password", h.resetPassword)

	// user superadmin endpoint
	superadminRouter.Get("/users", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "users"), h.getListUser)
	superadminRouter.Get("/users/:user_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "users"), h.getDetailUser)
	superadminRouter.Post("/users/create", h.middleware.UserBearer, h.middleware.RBACMiddleware("create", "users"), h.createUser)
	superadminRouter.Patch("/users/update/:user_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("update", "users"), h.updateUser)
	superadminRouter.Delete("/users/remove/:user_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("delete", "users"), h.removeUser)
}
