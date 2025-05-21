package rest

import (
	"strconv"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/expenses/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *expensesHandler) createExpenses(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateExpensesRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createExpenses - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createExpenses - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateExpenses(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("handler::createExpenses - Failed to create expenses")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *expensesHandler) getListExpenses(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
	)

	res, err := h.service.GetListExpenses(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListExpenses - Failed to get list expenses")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *expensesHandler) updateExpenses(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateExpensesRequest)
		ctx = c.Context()
	)

	expensesIDStr := c.Params("expenses_id")
	if expensesIDStr == "" {
		log.Warn().Msg("handler::updateExpenses - Expenses ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Expenses ID is required"))
	}

	id, err := strconv.Atoi(expensesIDStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::updateExpenses - Invalid expenses ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid expenses ID"))
	}

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateExpenses - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateExpenses - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.UpdateExpenses(ctx, req, id)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateExpenses - Failed to update expenses")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *expensesHandler) removeExpenses(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	expensesIDStr := c.Params("expenses_id")
	if expensesIDStr == "" {
		log.Warn().Msg("handler::removeExpenses - Expenses ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Expenses ID is required"))
	}

	id, err := strconv.Atoi(expensesIDStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::removeExpenses - Invalid expenses ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid expenses ID"))
	}

	err = h.service.RemoveExpenses(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("handler::removeExpenses - Failed to remove expenses")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}
