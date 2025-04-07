package rest

import "github.com/gofiber/fiber/v2"

func (h *productHandler) ProductRoute(userRouter, superadminRouter fiber.Router) {
	superadminRouter.Post("/product/create", h.middleware.UserBearer, h.middleware.RBACMiddleware("create", "product"), h.createProduct)
}
