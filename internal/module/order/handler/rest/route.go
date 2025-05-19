package rest

import "github.com/gofiber/fiber/v2"

func (h *orderHandler) OrderRoute(userRouter fiber.Router) {
	userRouter.Post("/order", h.middleware.UserBearer, h.createOrder)
	userRouter.Get("/order", h.middleware.UserBearer, h.getListOrder)
}
