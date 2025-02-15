package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/constants"
	externalNotification "github.com/Digitalkeun-Creative/be-dzikra-user-service/external/notification"
	redisPorts "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/infrastructure/redis"
	rolePorts "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/role/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/dto"
	user "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/entity"
	userPorts "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/ports"
	userProfile "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_profile/entity"
	userProfilePorts "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_profile/ports"
	userRole "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_role/entity"
	userRolePorts "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_role/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/utils"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ userPorts.UserService = &userService{}

type userService struct {
	db                    *sqlx.DB
	userRepository        userPorts.UserRepository
	roleRepository        rolePorts.RoleRepository
	userRoleRepository    userRolePorts.UserRoleRepository
	userProfileRepository userProfilePorts.UserProfileRepository
	redisRepository       redisPorts.RedisRepository
	externalNotification  externalNotification.ExternalNotification
}

func NewUserService(db *sqlx.DB, userRepository userPorts.UserRepository, roleRepository rolePorts.RoleRepository, userRoleRepository userRolePorts.UserRoleRepository, userProfileRepository userProfilePorts.UserProfileRepository, redisRepository redisPorts.RedisRepository, externalNotification externalNotification.ExternalNotification) *userService {
	return &userService{
		db:                    db,
		userRepository:        userRepository,
		roleRepository:        roleRepository,
		userRoleRepository:    userRoleRepository,
		userProfileRepository: userProfileRepository,
		redisRepository:       redisRepository,
		externalNotification:  externalNotification,
	}
}

func (s *userService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to hash password")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	req.Password = hashedPassword

	roleData, err := s.roleRepository.FindRoleByName(ctx, strings.ToUpper(constants.UserRole))
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to find role")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	userResult, err := s.userRepository.FindUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to find user by email")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	var res = new(user.User)

	tx, err := s.db.Begin()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::Register - Failed to rollback transaction")
			}
		}
	}()

	if userResult != nil {
		res.ID = userResult.ID
		res.FullName = userResult.FullName
		res.Username = userResult.Username
		res.Email = userResult.Email
	}

	if userResult == nil {
		log.Info().Any("payload", req).Msg("service::Register - User not registered")

		userID, err := utils.GenerateUUIDv7String()
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to generate UUID V7")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		userRoleID, err := utils.GenerateUUIDv7String()
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to generate UUID V7")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		userProfileID, err := utils.GenerateUUIDv7String()
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to generate UUID V7")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		userData := &user.User{
			ID:       userID,
			FullName: req.FullName,
			Username: req.Username,
			Password: req.Password,
			Email:    req.Email,
		}

		userRoleData := &userRole.UserRole{
			ID:     userRoleID,
			UserID: userID,
			RoleID: roleData.ID,
		}

		userProfileData := &userProfile.UserProfile{
			ID:          userProfileID,
			UserID:      userID,
			PhoneNumber: &req.PhoneNumber,
		}

		res, err = s.userRepository.InsertNewUser(ctx, tx, userData)
		if err != nil {
			if strings.Contains(err.Error(), constants.ErrUsernameAlreadyRegistered) {
				log.Error().Any("payload", req).Msg("service::Register - Username already registered")
				return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrUsernameAlreadyRegistered))
			}

			if strings.Contains(err.Error(), constants.ErrEmailAlreadyRegistered) {
				log.Error().Any("payload", req).Msg("service::Register - Email already registered")
				return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrEmailAlreadyRegistered))
			}

			log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to insert new user")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		err = s.userRoleRepository.InsertNewUserRole(ctx, tx, userRoleData)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to insert new user role")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		err = s.userProfileRepository.InsertNewUserProfile(ctx, tx, userProfileData)
		if err != nil {
			if strings.Contains(err.Error(), constants.ErrPhoneNumberAlreadyRegistered) {
				log.Error().Any("payload", req).Msg("service::Register - Phone number already registered")
				return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrPhoneNumberAlreadyRegistered))
			}

			log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to insert new user profile")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		err = s.redisRepository.Set(ctx, fmt.Sprintf("%s:%s", constants.RegisteredUserProgress, req.Username), true, 5*time.Minute)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to set data redis")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	} else if userResult.EmailVerifiedAt == nil && userResult.OtpNumberVerifiedAt == nil {
		log.Info().Any("payload", req).Msg("service::Register - Email not verified")

		key := fmt.Sprintf("%s:%s", constants.RegisteredUserProgress, req.Username)

		_, err = s.redisRepository.Get(ctx, key)
		if err != nil {
			if err == redis.Nil {
				log.Info().Any("payload", req).Msg("service::Register - OTP number is not found")
				err = s.redisRepository.Set(ctx, key, true, 5*time.Minute)
				if err != nil {
					log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to delete data redis")
					return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
				}
			} else {
				log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to get data redis")
				return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}
		}

		return nil, err_msg.NewCustomErrors(fiber.StatusUnprocessableEntity, err_msg.WithMessage(constants.ErrUserRegistrationInProgress))
	} else if userResult.EmailVerifiedAt != nil && userResult.OtpNumberVerifiedAt != nil {
		log.Error().Any("payload", req).Msg("service::Register - Email already verified")
		return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrEmailAlreadyRegistered))
	}

	otpNumber := utils.GenerateRandomOTP()

	otpRedisKey := fmt.Sprintf("%s:%s", constants.OTPNumberKey, req.Username)

	err = s.redisRepository.Set(ctx, otpRedisKey, otpNumber, 2*time.Minute)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to set data redis")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	go func() {
		err := s.externalNotification.SendNotification(context.Background(), req.Email, constants.RegisterTemplateName, map[string]string{
			"full_name":  req.FullName,
			"otp_number": otpNumber,
		})
		if err != nil {
			log.Error().Err(err).Msg("service::Register - Failed to send notification in background")
		}
	}()

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.RegisterResponse{
		Email:       res.Email,
		FullName:    res.FullName,
		PhoneNumber: req.PhoneNumber,
	}, nil
}
