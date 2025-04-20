package rest

import (
	"strconv"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *voucherHandler) createVoucher(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateVoucherRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createVoucher - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createVoucher - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateVoucher(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("handler::createVoucher - Failed to create voucher")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *voucherHandler) getListVoucher(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
	)

	res, err := h.service.GetListVoucher(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListVoucher - Failed to get list voucher")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *voucherHandler) updateVoucher(c *fiber.Ctx) error {
	ctx := c.Context()

	voucherIDStr := c.Params("voucher_id")
	if voucherIDStr == "" {
		log.Warn().Msg("handler::updateVoucher - Voucher ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Voucher ID is required"))
	}

	id, err := strconv.Atoi(voucherIDStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::updateVoucher - Invalid voucher ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid voucher ID"))
	}

	req := new(dto.CreateOrUpdateVoucherRequest)
	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateVoucher - parse body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateVoucher - validation")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.UpdateVoucher(ctx, id, req)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateVoucher - service error")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *voucherHandler) removeVoucher(c *fiber.Ctx) error {
	var (
		ctx          = c.Context()
		voucherIDStr = c.Params("voucher_id")
	)

	if strings.Contains(voucherIDStr, ":voucher_id") {
		log.Warn().Msg("handler::removeVoucher - invalid voucher ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid banner ID"))
	}

	id, err := strconv.Atoi(voucherIDStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::removeVoucher - invalid id parameter")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid id parameter"))
	}

	err = h.service.RemoveVoucher(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("handler::removeVoucher - failed to remove voucher")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}
