package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *provinceHandler) getListProvince(c *fiber.Ctx) error {
	res, err := h.service.GetListProvince(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("handler::getListProvince - Failed to get list province")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}
