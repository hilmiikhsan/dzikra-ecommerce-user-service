package rest

import (
	"strconv"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *subDistrict) getListSubDistrict(c *fiber.Ctx) error {
	var (
		ctx           = c.Context()
		districtIDStr = c.Params("district_id")
	)

	if strings.Contains(districtIDStr, ":district_id") {
		log.Warn().Msg("handler::getListSubDistrict - invalid district ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid district ID"))
	}

	districtID, err := strconv.Atoi(districtIDStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::getListSubDistrict - invalid id parameter")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid id parameter"))
	}

	res, err := h.service.GetListSubDistrict(ctx, districtID)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListSubDistrict - Failed to get list sub district")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}
