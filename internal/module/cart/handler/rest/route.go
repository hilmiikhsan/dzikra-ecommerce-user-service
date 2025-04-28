package rest

import "github.com/gofiber/fiber/v2"

func (h *cartHandler) CartRoute(userRouter fiber.Router) {
	userRouter.Post("/cart", h.middleware.UserBearer, h.addToCartItem)
	userRouter.Get("/cart", h.middleware.UserBearer, h.getListCart)
}
