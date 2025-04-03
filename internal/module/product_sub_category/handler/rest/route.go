package rest

import "github.com/gofiber/fiber/v2"

func (h *productSubCategoryHandler) ProductSubCategoryRoute(userRouter, superadminRouter fiber.Router) {
	superadminRouter.Post("/product_category/:category_id/product_subcategory/create", h.middleware.UserBearer, h.middleware.RBACMiddleware("create", "subcategory"), h.createProductSubCategory)
}
