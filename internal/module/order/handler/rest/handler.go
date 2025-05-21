package rest

import (
	"strconv"
	"strings"

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
		log.Warn().Err(err).Msg("handler::createOrder - Invalid request body")
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

func (h *orderHandler) getWaybill(c *fiber.Ctx) error {
	var (
		ctx     = c.Context()
		orderID = c.Params("order_id")
	)

	if strings.Contains(orderID, ":order_id") {
		log.Warn().Msg("handler::getWaybill - invalid order ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid order ID"))
	}

	res, err := h.service.GetWaybillDetails(ctx, orderID)
	if err != nil {
		log.Error().Err(err).Msg("handler::getWaybill - Failed to get waybill")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *orderHandler) getListOrderTransaction(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
		status = c.Query("status", "")
	)

	res, err := h.service.GetListOrderTransaction(ctx, page, limit, search, status)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListOrderTransaction - Failed to get list order transaction")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *orderHandler) getOrderWaybillTransaction(c *fiber.Ctx) error {
	var (
		ctx     = c.Context()
		orderID = c.Params("order_id")
	)

	if strings.Contains(orderID, ":order_id") {
		log.Warn().Msg("handler::getWaybill - invalid order ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid order ID"))
	}

	res, err := h.service.GetWaybillDetails(ctx, orderID)
	if err != nil {
		log.Error().Err(err).Msg("handler::getWaybill - Failed to get waybill")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *orderHandler) updateOrderShippingNumber(c *fiber.Ctx) error {
	var (
		ctx     = c.Context()
		orderID = c.Params("order_id")
		req     = new(dto.UpdateOrderShippingNumberRequest)
	)

	if strings.Contains(orderID, ":order_id") {
		log.Warn().Msg("handler::updateOrderShippingNumber - invalid order ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid order ID"))
	}

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateOrderShippingNumber - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateOrderShippingNumber - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.UpdateOrderShippingNumber(ctx, req, orderID)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateOrderShippingNumber - Failed to update order shipping number")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *orderHandler) updateOrderStatusTransaction(c *fiber.Ctx) error {
	var (
		ctx     = c.Context()
		orderID = c.Params("order_id")
		req     = new(dto.UpdateOrderStatusTransactionRequest)
		locals  = middleware.GetLocals(c)
	)

	if strings.Contains(orderID, ":order_id") {
		log.Warn().Msg("handler::updateOrderStatusTransaction - invalid order ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid order ID"))
	}

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateOrderStatusTransaction - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateOrderStatusTransaction - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err := h.service.UpdateOrderStatusTransaction(ctx, req, orderID, locals.FullName, locals.Email)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateOrderStatusTransaction - Failed to update order status transaction")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}
