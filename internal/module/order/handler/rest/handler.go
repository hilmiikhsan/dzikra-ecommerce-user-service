package rest

import (
	"strconv"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/order/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *orderHandler) createOrder(c *fiber.Ctx) error {
	var (
		req    = new(dto.CreateOrderRequest)
		ctx    = c.Context()
		locals = middleware.GetLocals(c)
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createOrder - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createAddress - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	addressID, err := strconv.Atoi(req.AddressID)
	if err != nil {
		log.Warn().Msg("handler::createOrder - Invalid address ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid address ID"))
	}

	var voucherID int
	if req.VoucherID != nil {
		voucherID, err = strconv.Atoi(*req.VoucherID)
		if err != nil {
			log.Warn().Msg("handler::createOrder - Invalid voucher ID")
			return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid voucher ID"))
		}
	}

	res, err := h.service.CreateOrder(ctx, req, locals, addressID, voucherID)
	if err != nil {
		log.Error().Err(err).Msg("handler::createOrder - Failed to create order")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *orderHandler) getListOrder(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
		status = c.Query("status", "")
		locals = middleware.GetLocals(c)
	)

	res, err := h.service.GetListOrder(ctx, page, limit, search, status, locals.UserID)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListOrder - Failed to get list order")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}
