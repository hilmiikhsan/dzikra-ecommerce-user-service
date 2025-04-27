package rest

import "github.com/gofiber/fiber/v2"

func (h *addressHandler) AddressRoute(userRouter fiber.Router) {
	userRouter.Post("/address/create", h.middleware.UserBearer, h.middleware.RBACMiddleware("create", "address"), h.createAddress)
	userRouter.Patch("/address/update/:address_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("update", "address"), h.updateAddress)
}
