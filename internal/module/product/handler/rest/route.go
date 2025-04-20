package rest

import "github.com/gofiber/fiber/v2"

func (h *productHandler) ProductRoute(publicRouter, superadminRouter fiber.Router) {
	// superadmin endpoint
	superadminRouter.Post("/product/create", h.middleware.UserBearer, h.middleware.RBACMiddleware("create", "product"), h.createProduct)
	superadminRouter.Patch("/product/update/:product_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("update", "product"), h.updateProduct)
	superadminRouter.Get("/product", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "product"), h.getListProduct)
	superadminRouter.Get("/product/:product_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "product"), h.getDetailProduct)
	superadminRouter.Delete("/product/remove/:product_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("delete", "product"), h.removeProduct)

	// public endpoint
	publicRouter.Get("/product", h.getListProduct)
	publicRouter.Get("/product/:product_id", h.getDetailProduct)
}
