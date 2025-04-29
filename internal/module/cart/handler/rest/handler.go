package rest

import (
	"strconv"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *cartHandler) addToCartItem(c *fiber.Ctx) error {
	var (
		req = new(dto.AddOrUpdateCartItemRequest)
		ctx = c.Context()
	)

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		log.Error().Msg("middleware::addToCartItem - user_id not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(constants.ErrAccessTokenIsRequired))
	}

	req.UserID = userID

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::addToCartItem - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::addToCartItem - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.AddCartItem(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("handler::addToCartItem - Failed to add item to cart")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *cartHandler) getListCart(c *fiber.Ctx) error {
	ctx := c.Context()

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		log.Error().Msg("middleware::getListCart - user_id not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(constants.ErrAccessTokenIsRequired))
	}

	res, err := h.service.GetListCart(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListCart - Failed to get list cart")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *cartHandler) updateCartItem(c *fiber.Ctx) error {
	var (
		req       = new(dto.AddOrUpdateCartItemRequest)
		ctx       = c.Context()
		cartIDStr = c.Params("cart_id")
	)

	if strings.Contains(cartIDStr, ":cart_id") {
		log.Warn().Msg("handler::updateCartItem - invalid cart ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid cart ID"))
	}

	id, err := strconv.Atoi(cartIDStr)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateCartItem - Failed to convert id to int")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to convert id to int"))
	}

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		log.Error().Msg("middleware::updateCartItem - user_id not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(constants.ErrAccessTokenIsRequired))
	}

	req.UserID = userID

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateCartItem - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateCartItem - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.UpdateCartItem(ctx, req, id)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateCartItem - Failed to update cart item")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}
