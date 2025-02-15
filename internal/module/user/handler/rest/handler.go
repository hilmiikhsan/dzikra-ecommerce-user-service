package rest

import (
	externalNotification "github.com/Digitalkeun-Creative/be-dzikra-user-service/external/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/adapter"
	redisRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/infrastructure/redis"
	roleRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/ports"
	userRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/service"
	userProfileRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_profile/repository"
	userRoleRepository "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_role/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type userHandler struct {
	service ports.UserService
	// middleware middleware.AuthMiddleware
	validator adapter.Validator
}

func NewUserHandler() *userHandler {
	var handler = new(userHandler)

	// validator
	validator := adapter.Adapters.Validator

	// external
	externalNotification := &externalNotification.External{}

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.DzikraRedis)

	// jwt
	// jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	// middlewareHandler := middleware.NewAuthMiddleware(jwt)

	// repository
	userRepository := userRepository.NewUserRepository(adapter.Adapters.DzikraPostgres)
	roleRepository := roleRepository.NewRoleRepository(adapter.Adapters.DzikraPostgres)
	userRoleRepository := userRoleRepository.NewUserRoleRepository(adapter.Adapters.DzikraPostgres)
	userProfileRepository := userProfileRepository.NewUserProfileRepository(adapter.Adapters.DzikraPostgres)

	// service
	userService := service.NewUserService(
		adapter.Adapters.DzikraPostgres,
		userRepository,
		roleRepository,
		userRoleRepository,
		userProfileRepository,
		redisRepository,
		externalNotification,
	)

	// handler
	handler.service = userService
	// handler.middleware = *middlewareHandler
	handler.validator = validator

	return handler
}

func (h *userHandler) UserRoute(router fiber.Router) {
	router.Post("/register", h.register)
}

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
