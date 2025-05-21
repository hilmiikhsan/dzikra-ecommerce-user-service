package rest

import "github.com/gofiber/fiber/v2"

func (h *orderHandler) OrderRoute(userRouter, superadminRouter fiber.Router) {
	// user
	userRouter.Post("/order", h.middleware.UserBearer, h.createOrder)
	userRouter.Get("/order", h.middleware.UserBearer, h.getListOrder)
	userRouter.Get("/order/:order_id/waybill", h.middleware.UserBearer, h.getWaybill)

	// superadmin
	superadminRouter.Get("/order", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "order"), h.getListOrderTransaction)
	superadminRouter.Get("/order/:order_id/waybill", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "order"), h.getOrderWaybillTransaction)
	superadminRouter.Put("/order/:order_id/shipping-number", h.middleware.UserBearer, h.middleware.RBACMiddleware("update", "order"), h.updateOrderShippingNumber)
	superadminRouter.Put("/order/:order_id/status", h.middleware.UserBearer, h.middleware.RBACMiddleware("update", "order"), h.updateOrderStatusTransaction)
}
