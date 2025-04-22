package rest

import "github.com/gofiber/fiber/v2"

func (h *provinceHandler) ProvinceRoute(publicRouter fiber.Router) {
	publicRouter.Get("/location/province", h.getListProvince)
}
