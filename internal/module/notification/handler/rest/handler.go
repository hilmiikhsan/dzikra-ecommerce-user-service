package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/notification/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *notificationHandler) createNotification(c *fiber.Ctx) error {
	var (
		req    = new(dto.CreateNotificationRequest)
		ctx    = c.Context()
		locals = middleware.GetLocals(c)
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createNotification - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.UserID = locals.UserID

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createNotification - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err := h.service.CreateNotification(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("handler::createNotification - Failed to create notification")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}

func (h *notificationHandler) getListNotification(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
	)

	res, err := h.service.GetListNotification(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListNotification - Failed to get list notification")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *notificationHandler) sendFcmBatchNotification(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req = new(dto.SendFcmBatchRequest)
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::sendFcmBatchNotification - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::sendFcmBatchNotification - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err := h.service.SendFcmBatchNotification(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("handler::sendFcmBatchNotification - Failed to send FCM batch notification")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}
