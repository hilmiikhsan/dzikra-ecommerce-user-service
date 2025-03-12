package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *superAdminHandler) createRolePermission(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateRolePermissionRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createRolePermission - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createRolePermission - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateRolePermission(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::createRolePermission - Failed to create role permission")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *superAdminHandler) getListRole(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	res, err := h.service.GetListRole(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListRole - Failed to get list role")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *superAdminHandler) getListApplication(c *fiber.Ctx) error {
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
