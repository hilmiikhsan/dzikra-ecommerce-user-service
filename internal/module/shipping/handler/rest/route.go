package rest

import "github.com/gofiber/fiber/v2"

func (h *shippingHandler) ShippingRoute(userRouter fiber.Router) {
	userRouter.Post("/shipping/cost", h.middleware.UserBearer, h.calculateShippingCost)
}
