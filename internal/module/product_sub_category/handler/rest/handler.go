package rest

import (
	"strconv"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *productSubCategoryHandler) createProductSubCategory(c *fiber.Ctx) error {
	var (
		req   = new(dto.CreateOrUpdateProductSubCategoryRequest)
		ctx   = c.Context()
		idStr = c.Params("category_id")
	)

	if strings.Contains(idStr, ":category_id") {
		log.Warn().Msg("handler::createProductSubCategory - Invalid product category ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid product category ID"))
	}

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createProductSubCategory - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createProductSubCategory - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	categoryID, _ := strconv.Atoi(c.Params("category_id"))

	res, err := h.service.CreateProductSubCategory(ctx, req, categoryID)
	if err != nil {
		log.Error().Err(err).Msg("handler::createProductSubCategory - Failed to create product sub category")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}
