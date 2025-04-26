package rest

import "github.com/gofiber/fiber/v2"

func (h *subDistrict) SubDistrictRoute(publicRouter fiber.Router) {
	publicRouter.Get("/location/subdistrict/:district_id", h.getListSubDistrict)
}
