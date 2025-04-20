package rest

import "github.com/gofiber/fiber/v2"

func (h *productSubCategoryHandler) ProductSubCategoryRoute(publicRouter, superadminRouter fiber.Router) {
	// superadmin endpoint
	superadminRouter.Post("/product_category/:category_id/product_subcategory/create", h.middleware.UserBearer, h.middleware.RBACMiddleware("create", "subcategory"), h.createProductSubCategory)
	superadminRouter.Patch("/product_category/:category_id/product_subcategory/update/:subcategory_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("update", "subcategory"), h.updateProductSubCategory)
	superadminRouter.Get("/product_category/:category_id/product_subcategory", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "subcategory"), h.getListProductSubCategory)
	superadminRouter.Delete("/product_subcategory/remove/:subcategory_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("delete", "subcategory"), h.removeProductSubCategory)

	// public endpoint
	publicRouter.Get("/product_category/:category_id/product_subcategory", h.getListProductSubCategory)
}
