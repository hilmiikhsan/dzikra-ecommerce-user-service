package rest

import "github.com/gofiber/fiber/v2"

func (h *bannerHandler) BannerRoute(userRouter, superadminRouter fiber.Router) {
	superadminRouter.Post("/banner/create", h.middleware.UserBearer, h.middleware.RBACMiddleware("create", "banner"), h.createBanner)
	superadminRouter.Get("/banner", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "banner"), h.getListBanner)
	superadminRouter.Patch("/banner/update/:banner_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("update", "banner"), h.updateBanner)
}
