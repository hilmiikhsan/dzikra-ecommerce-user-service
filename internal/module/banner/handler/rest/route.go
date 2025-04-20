package rest

import "github.com/gofiber/fiber/v2"

func (h *bannerHandler) BannerRoute(publicRouter, superadminRouter fiber.Router) {
	// superadmin endpoint
	superadminRouter.Post("/banner/create", h.middleware.UserBearer, h.middleware.RBACMiddleware("create", "banner"), h.createBanner)
	superadminRouter.Get("/banner", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "banner"), h.getListBanner)
	superadminRouter.Patch("/banner/update/:banner_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("update", "banner"), h.updateBanner)
	superadminRouter.Delete("/banner/remove/:banner_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("delete", "banner"), h.removeBanner)

	// public endpoint
	publicRouter.Get("/banner", h.getListBanner)
}
