package rest

import "github.com/gofiber/fiber/v2"

func (h *expensesHandler) ExpensesRoute(superadminRouter fiber.Router) {
	superadminRouter.Post("/expenses/create", h.middleware.UserBearer, h.middleware.RBACMiddleware("create", "expenses"), h.createExpenses)
	superadminRouter.Get("/expenses", h.middleware.UserBearer, h.middleware.RBACMiddleware("read", "expenses"), h.getListExpenses)
	superadminRouter.Patch("/expenses/update/:expenses_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("update", "expenses"), h.updateExpenses)
	superadminRouter.Delete("/expenses/remove/:expenses_id", h.middleware.UserBearer, h.middleware.RBACMiddleware("delete", "expenses"), h.removeExpenses)
}
