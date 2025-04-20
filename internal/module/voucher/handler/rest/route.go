package rest

import "github.com/gofiber/fiber/v2"

func (h *voucherHandler) VoucherRoute(userRouter, superadminRouter fiber.Router) {
	superadminRouter.Post("/voucher/create", h.middleware.UserBearer, h.middleware.RBACMiddleware("create", "voucher"), h.createVoucher)
	superadminRouter.Get("/voucher", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "voucher"), h.getListVoucher)
	superadminRouter.Patch("/voucher/update/:voucher_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("update", "voucher"), h.updateVoucher)
	superadminRouter.Delete("/voucher/remove/:voucher_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("delete", "voucher"), h.removeVoucher)
}
