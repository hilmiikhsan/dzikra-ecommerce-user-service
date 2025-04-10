package rest

import "github.com/gofiber/fiber/v2"

func (h *productCategoryHandler) ProductCategoryRoute(userRouter, superadminRouter fiber.Router) {
	superadminRouter.Get("/product_category", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "category"), h.getListProductCategory)
	superadminRouter.Get("/category/:id_category", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "category"), h.getDetailProductCategory)
	superadminRouter.Post("/product_category/create", h.middleware.UserBearer, h.middleware.RBACMiddleware("create", "category"), h.createProductCategory)
	superadminRouter.Patch("/product_category/update/:product_category_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("update", "category"), h.updateProductCategory)
	superadminRouter.Delete("/product_category/remove/:product_category_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("delete", "category"), h.removeProductCategory)
}
