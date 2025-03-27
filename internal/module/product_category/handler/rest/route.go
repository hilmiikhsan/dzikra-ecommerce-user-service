package rest

import "github.com/gofiber/fiber/v2"

func (h *productCategoryHandler) ProductCategoryRoute(userRouter, superadminRouter fiber.Router) {
	superadminRouter.Get("/product_category", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "category"), h.getListProductCategory)
	superadminRouter.Post("/product_category/create", h.middleware.UserBearer, h.middleware.RBACMiddleware("create", "category"), h.createProductCategory)
	superadminRouter.Patch("/product_category/update/:product_category_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("update", "category"), h.updateProductCategory)
}
