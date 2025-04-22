package rest

import "github.com/gofiber/fiber/v2"

func (h *cityHandler) CityRoute(publicRouter fiber.Router) {
	publicRouter.Get("/location/city/:province_id", h.getListCity)
}
