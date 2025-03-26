package rest

import (
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/role/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *roleHandler) getDetailRole(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		roleID = c.Params("role_id")
	)

	if strings.Contains(roleID, ":role_id") {
		log.Warn().Msg("handler::getDetailRole - Invalid role ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid role ID"))
	}

	res, err := h.service.GetDetailRole(ctx, roleID)
	if err != nil {
		log.Error().Err(err).Any("roleID", roleID).Msg("handler::getDetailRole - Failed to get detail role")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *roleHandler) getListRole(c *fiber.Ctx) error {
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

func (h *roleHandler) createRolePermission(c *fiber.Ctx) error {
	var (
		req = new(dto.RolePermissionRequest)
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

func (h *roleHandler) updateRolePermission(c *fiber.Ctx) error {
	var (
		req    = new(dto.SoftDeleteRolePermissionRequest)
		ctx    = c.Context()
		roleID = c.Params("role_id")
	)

	if strings.Contains(roleID, ":role_id") {
		log.Warn().Msg("handler::updateRolePermission - Invalid role ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid role ID"))
	}

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateRolePermission - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateRolePermission - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.UpdateRolePermission(ctx, req, roleID)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::updateRolePermission - Failed to update role permission")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *roleHandler) removeRolePermission(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		roleID = c.Params("role_id")
	)

	if strings.Contains(roleID, ":role_id") {
		log.Warn().Msg("handler::removeRolePermission - Invalid role ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid role ID"))
	}

	err := h.service.RemoveRolePermission(ctx, roleID)
	if err != nil {
		log.Error().Err(err).Any("roleID", roleID).Msg("handler::removeRolePermission - Failed to remove role permission")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}
