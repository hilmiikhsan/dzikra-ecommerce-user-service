package rest

import (
	"strconv"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *cityHandler) getListCity(c *fiber.Ctx) error {
	var (
		ctx           = c.Context()
		provinceIDStr = c.Params("province_id")
	)

	if strings.Contains(provinceIDStr, ":province_id") {
		log.Warn().Msg("handler::getListCity - invalid province ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid province ID"))
	}

	provinceID, err := strconv.Atoi(provinceIDStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::getListCity - invalid id parameter")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid id parameter"))
	}

	res, err := h.service.GetListCity(ctx, provinceID)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListCity - Failed to get list city")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}
