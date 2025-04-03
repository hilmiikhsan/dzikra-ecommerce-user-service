package rest

import (
	"strconv"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *productCategoryHandler) getListProductCategory(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	res, err := h.service.GetListProductCategory(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListProductCategory - Failed to get list product category")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *productCategoryHandler) createProductCategory(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateProductCategoryRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createProductCategory - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createProductCategory - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateProductCategory(ctx, req.Category)
	if err != nil {
		log.Error().Err(err).Msg("handler::createProductCategory - Failed to create product category")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *productCategoryHandler) updateProductCategory(c *fiber.Ctx) error {
	var (
		req   = new(dto.CreateOrUpdateProductCategoryRequest)
		ctx   = c.Context()
		idStr = c.Params("product_category_id")
	)

	if strings.Contains(idStr, ":product_category_id") {
		log.Warn().Msg("handler::updateProductCategory - Invalid product category ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid product category ID"))
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::updateProductCategory - Invalid id parameter")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid id parameter"))
	}

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateProductCategory - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateProductCategory - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.UpdateProductCategory(ctx, id, req.Category)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateProductCategory - Failed to update product category")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *productCategoryHandler) removeProductCategory(c *fiber.Ctx) error {
	var (
		ctx   = c.Context()
		idStr = c.Params("product_category_id")
	)

	if strings.Contains(idStr, ":product_category_id") {
		log.Warn().Msg("handler::removeProductCategory - Invalid product category ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid product category ID"))
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::removeProductCategory - Invalid id parameter")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid id parameter"))
	}

	err = h.service.RemoveProductCategory(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("handler::removeProductCategory - Failed to remove product category")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}
