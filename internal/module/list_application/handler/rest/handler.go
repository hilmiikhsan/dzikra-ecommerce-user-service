package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *applicationHandler) getListApplication(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	res, err := h.service.GetListApplication(ctx)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListApplication - Failed to get list application")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *applicationHandler) getListPermissionByApp(c *fiber.Ctx) error {
	var (
		ctx         = c.Context()
		appIDsParam = c.Query("appid")
	)

	res, err := h.service.GetListPermissionByApp(ctx, appIDsParam)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListPermissionByApp - Failed to get list permission by app")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}
