package rest

import "github.com/gofiber/fiber/v2"

func (h *notificationHandler) NotificationRoute(superadminRouter fiber.Router) {
	superadminRouter.Post("/notification/create", h.middleware.UserBearer, h.middleware.RBACMiddleware("create", "notification"), h.createNotification)
	superadminRouter.Get("/notification", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "notification"), h.getListNotification)
	superadminRouter.Post("/notification/send_fcm_batch", h.middleware.UserBearer, h.sendFcmBatchNotification)
}
