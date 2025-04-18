package rest

import (
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/user/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *userHandler) register(c *fiber.Ctx) error {
	var (
		req = new(dto.RegisterRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::register - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::register - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.Register(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::register - Failed to register user")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *userHandler) verification(c *fiber.Ctx) error {
	var (
		req = new(dto.VerificationRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::verification - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::verification - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.Verification(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::verification - Failed to verify user")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *userHandler) sendOtpNumberVerification(c *fiber.Ctx) error {
	var (
		req = new(dto.SendOtpNumberVerificationRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::sendOtpNumberVerification - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::sendOtpNumberVerification - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.SendOtpNumberVerification(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::sendOtpNumberVerification - Failed to send OTP number verification")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *userHandler) login(c *fiber.Ctx) error {
	var (
		req = new(dto.LoginRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::login - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::login - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.Login(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::login - Failed to login user")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *userHandler) logout(c *fiber.Ctx) error {
	var (
		ctx         = c.Context()
		accessToken = c.Get(constants.HeaderAuthorization)
		locals      = middleware.GetLocals(c)
	)

	if len(accessToken) > 7 {
		accessToken = accessToken[7:]
	}

	err := h.service.Logout(ctx, accessToken, locals)
	if err != nil {
		log.Error().Err(err).Any("access_token", accessToken).Msg("handler::logout - Failed to logout user")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}

func (h *userHandler) getCurrentUser(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		locals = middleware.GetLocals(c)
	)

	res, err := h.service.GetCurrentUser(ctx, locals)
	if err != nil {
		log.Error().Err(err).Msg("handler::getCurrentUser - Failed to get current user")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *userHandler) refreshToken(c *fiber.Ctx) error {
	var (
		ctx         = c.Context()
		accessToken = c.Get(constants.HeaderAuthorization)
		locals      = middleware.GetLocals(c)
	)

	if len(accessToken) > 7 {
		accessToken = accessToken[7:]
	}

	res, err := h.service.RefreshToken(ctx, accessToken, locals)
	if err != nil {
		log.Error().Err(err).Any("access_token", accessToken).Msg("handler::refreshToken - Failed to refresh token")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *userHandler) forgotPassword(c *fiber.Ctx) error {
	var (
		req = new(dto.SendOtpNumberVerificationRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::forgotPassword - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::forgotPassword - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.ForgotPassword(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::forgotPassword - Failed to forgot password")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *userHandler) resetPassword(c *fiber.Ctx) error {
	var (
		req = new(dto.ResetPasswordRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::resetPassword - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::resetPassword - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err := h.service.ResetPassword(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::resetPassword - Failed to reset password")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}

func (h *userHandler) getListUser(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
	)

	res, err := h.service.GetListUser(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::GetListUser - Failed to get list user")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *userHandler) getDetailUser(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		userID = c.Params("user_id")
	)

	if strings.Contains(userID, ":user_id") {
		log.Warn().Msg("handler::getDetailRole - Invalid user ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid user ID"))
	}

	res, err := h.service.GetDetailUser(ctx, userID)
	if err != nil {
		log.Error().Err(err).Any("user_id", userID).Msg("handler::getDetailUser - Failed to get detail user")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *userHandler) createUser(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateUserRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createUser - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createUser - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateUser(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::createUser - Failed to create user")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *userHandler) updateUser(c *fiber.Ctx) error {
	var (
		req    = new(dto.CreateOrUpdateUserRequest)
		ctx    = c.Context()
		userID = c.Params("user_id")
	)

	if strings.Contains(userID, ":user_id") {
		log.Warn().Msg("handler::updateUser - Invalid user ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid user ID"))
	}

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateUser - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateUser - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.UpdateUser(ctx, userID, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::updateUser - Failed to update user")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *userHandler) removeUser(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		userID = c.Params("user_id")
		locals = middleware.GetLocals(c)
	)

	if strings.Contains(userID, ":user_id") {
		log.Warn().Msg("handler::removeUser - Invalid user ID")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid user ID"))
	}

	if locals.UserID == userID {
		log.Warn().Any("user_id", userID).Msg("handler::removeUser - Cannot remove own user")
		return c.Status(fiber.StatusForbidden).JSON(response.Error("Cannot remove own user"))
	}

	err := h.service.RemoveUser(ctx, userID)
	if err != nil {
		log.Error().Err(err).Any("user_id", userID).Msg("handler::removeUser - Failed to remove user")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}
