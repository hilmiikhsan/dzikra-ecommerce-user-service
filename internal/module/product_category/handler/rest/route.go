package rest

import "github.com/gofiber/fiber/v2"

func (h *productCategoryHandler) ProductCategoryRoute(userRouter, superadminRouter fiber.Router) {
	superadminRouter.Get("/product_category", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "category"), h.getListProductCategory)
}
