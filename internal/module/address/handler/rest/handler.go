package rest

import (
	"strconv"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *addressHandler) createAddress(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateAddressRequest)
		ctx = c.Context()
	)

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		log.Error().Msg("middleware::createAddress - user_id not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(constants.ErrAccessTokenIsRequired))
	}

	req.UserID = userID

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createAddress - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createAddress - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	_, err := strconv.Atoi(req.ProvinceID)
	if err != nil {
		log.Warn().Msg("handler::createAddress - Invalid province ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid province ID"))
	}

	_, err = strconv.Atoi(req.CityID)
	if err != nil {
		log.Warn().Msg("handler::createAddress - Invalid city ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid city ID"))
	}

	_, err = strconv.Atoi(req.SubDistrictID)
	if err != nil {
		log.Warn().Msg("handler::createAddress - Invalid sub district ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid sub district ID"))
	}

	res, err := h.service.CreateAddress(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("handler::createAddress - Failed to create address")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *addressHandler) updateAddress(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateAddressRequest)
		ctx = c.Context()
	)

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		log.Error().Msg("middleware::updateAddress - user_id not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(constants.ErrAccessTokenIsRequired))
	}

	req.UserID = userID

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateAddress - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateAddress - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	addressID, err := strconv.Atoi(c.Params("address_id"))
	if err != nil {
		log.Warn().Msg("handler::updateAddress - Invalid address ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid address ID"))
	}

	res, err := h.service.UpdateAddress(ctx, req, addressID)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateAddress - Failed to update address")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(res, ""))
}

func (h *addressHandler) removeAddress(c *fiber.Ctx) error {
	var ctx = c.Context()

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		log.Error().Msg("middleware::removeAddress - user_id not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(constants.ErrAccessTokenIsRequired))
	}

	addressID, err := strconv.Atoi(c.Params("address_id"))
	if err != nil {
		log.Warn().Msg("handler::removeAddress - Invalid address ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid address ID"))
	}

	err = h.service.RemoveAddress(ctx, addressID, userID)
	if err != nil {
		log.Error().Err(err).Msg("handler::removeAddress - Failed to remove address")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}
