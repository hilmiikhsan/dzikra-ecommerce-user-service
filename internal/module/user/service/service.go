package service

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/middleware"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/dto"
	user "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user/entity"
	userFcmToken "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_fcm_token/entity"
	userProfile "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_profile/entity"
	userRole "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/module/user_role/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/jwt_handler"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/utils"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

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

	userResult, err := s.userRepository.FindByEmail(ctx, req.Email)
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
			} else {
				log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to get data redis")
				return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}
		} else {
			return nil, err_msg.NewCustomErrors(fiber.StatusUnprocessableEntity, err_msg.WithMessage(constants.ErrUserRegistrationInProgress))
		}
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

func (s *userService) Verification(ctx context.Context, req *dto.VerificationRequest) (*dto.VerificationResponse, error) {
	userData, err := s.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Verification - Failed to find user by email")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage("Failed to find user by email"))
	}

	if userData == nil {
		log.Error().Any("payload", req).Msg("service::Verification - User not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrEmailOrOTPNumberIsIncorrect))
	}

	key := fmt.Sprintf("%s:%s", constants.OTPNumberKey, userData.Username)
	blockedKey := fmt.Sprintf("%s:%s", constants.OTPBlockedKey, userData.Username)
	attemptKey := fmt.Sprintf("%s:%s", constants.OTPAttemptKey, userData.Username)

	blockTTL, err := s.redisRepository.TTL(ctx, blockedKey)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Verification - Failed to check blocked status in Redis")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	if blockTTL > 0 {
		return nil, err_msg.NewCustomErrors(fiber.StatusTooManyRequests, err_msg.WithMessage(fmt.Sprintf("You are blocked. Please try again in %d seconds.", int(blockTTL.Seconds()))))
	}

	otpNumber, err := s.redisRepository.Get(ctx, key)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Verification - Failed to get data from Redis")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrOTPNumberIsAlreadyExpired))
	}

	if otpNumber != req.Otp {
		attemptCount, err := s.redisRepository.Incr(ctx, attemptKey)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("service::Verification - Failed to increment attempt count")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		_ = s.redisRepository.Expire(ctx, attemptKey, 10*time.Minute)

		if attemptCount >= 3 {
			var blockDuration time.Duration
			switch attemptCount {
			case 3:
				blockDuration = 1 * time.Minute
			case 6:
				blockDuration = 3 * time.Minute
			default:
				blockDuration = 5 * time.Minute
			}

			err = s.redisRepository.Set(ctx, blockedKey, "blocked", blockDuration)
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("service::Verification - Failed to set blocked status in Redis")
				return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}

			return nil, err_msg.NewCustomErrors(fiber.StatusTooManyRequests, err_msg.WithMessage(fmt.Sprintf("You are blocked. Please try again in %d seconds.", int(blockDuration.Seconds()))))
		}

		return nil, err_msg.NewCustomErrors(fiber.StatusUnprocessableEntity, err_msg.WithMessage(constants.ErrInvalidOtpNumber))
	}

	emailVerifiedAt, err := s.userRepository.UpdateVerificationUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Verification - Failed to update verification user by email")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	keysToDelete := []string{key, blockedKey, attemptKey}
	for _, k := range keysToDelete {
		err := s.redisRepository.Del(ctx, k)
		if err != nil {
			log.Error().Err(err).Msgf("service::Verification - Failed to delete Redis key: %s", k)
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	log.Info().Any("user", userData).Msg("Verification successful and all Redis keys deleted")

	formattedVerifiedAt := utils.FormatToWIB(emailVerifiedAt)

	return &dto.VerificationResponse{
		Email: req.Email,
		EmailConfirmed: dto.EmailConfirmed{
			IsConfirm: true,
			CreatedAt: formattedVerifiedAt,
		},
	}, nil
}

func (s *userService) SendOtpNumberVerification(ctx context.Context, req *dto.SendOtpNumberVerificationRequest) (*dto.SendOtpNumberVerificationResponse, error) {
	userData, err := s.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::SendOtpNumberVerification - Failed to find user by email")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage("Failed to find user by email"))
	}

	if userData == nil {
		log.Error().Any("payload", req).Msg("service::SendOtpNumberVerification - User not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnprocessableEntity, err_msg.WithMessage(constants.ErrEmailNotRegistered))
	} else if userData.EmailVerifiedAt != nil && userData.OtpNumberVerifiedAt != nil {
		log.Error().Any("payload", req).Msg("service::SendOtpNumberVerification - User already verified")
		return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrUserAlreadyVerified))
	}

	key := fmt.Sprintf("%s:%s", constants.OTPNumberKey, userData.Username)

	_, err = s.redisRepository.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			generateNewOtpNumber := utils.GenerateRandomOTP()

			err := s.redisRepository.Set(ctx, key, generateNewOtpNumber, 2*time.Minute)
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("service::SendOtpNumberVerification - Failed to set data redis")
				return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}

			go func() {
				err := s.externalNotification.SendNotification(context.Background(), req.Email, constants.RegisterTemplateName, map[string]string{
					"full_name":  userData.FullName,
					"otp_number": generateNewOtpNumber,
				})
				if err != nil {
					log.Error().Err(err).Msg("service::SendOtpNumberVerification - Failed to send notification in background")
				}
			}()

			return &dto.SendOtpNumberVerificationResponse{
				Otp: "OK",
			}, nil
		}

		log.Error().Err(err).Any("payload", req).Msg("service::SendOtpNumberVerification - Failed to get data redis")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	log.Warn().Any("payload", req).Msg("service::SendOtpNumberVerification - Too many requests for the same email, please wait until the current OTP expires")
	return nil, err_msg.NewCustomErrors(fiber.StatusTooManyRequests, err_msg.WithMessage(constants.ErrTooManyReuqestOTPNumber))
}

func (s *userService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthUserResponse, error) {
	// define variable
	var (
		res = new(dto.AuthUserResponse)
	)

	// find user by email
	userResult, err := s.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrEmailOrPasswordIsIncorrect) {
			log.Error().Any("payload", req).Msg("service::Login - Email is incorrect")
			return nil, err_msg.NewCustomErrors(fiber.StatusUnprocessableEntity, err_msg.WithMessage(constants.ErrEmailOrPasswordIsIncorrect))
		}

		log.Error().Err(err).Any("payload", req).Msg("service::Login - Failed to find user by email")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check password
	if !utils.ComparePassword(userResult.Password, req.Password) {
		log.Error().Any("payload", req).Msg("service::Login - Password is incorrect")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnprocessableEntity, err_msg.WithMessage(constants.ErrEmailOrPasswordIsIncorrect))
	}

	// generate token
	result, err := s.jwt.GenerateTokenString(ctx, jwt_handler.CostumClaimsPayload{
		UserID:     userResult.ID.String(),
		Username:   userResult.Username,
		Email:      userResult.Email,
		FullName:   userResult.FullName,
		SessionID:  utils.GenerateSessionUUID(),
		DeviceID:   req.DeviceID,
		DeviceType: req.DeviceType,
		FcmToken:   req.FcmToken,
	})
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Login - Failed to generate token string")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get user profile
	userProfileResult, err := s.userProfileRepository.FindByUserID(ctx, userResult.ID.String())
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Login - Failed to find user profile by user id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get user role ids
	userRoleIDs, err := s.userRoleRepository.FindByUserID(ctx, userResult.ID.String())
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Login - Failed to find user role by user id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get user role permission
	userRolePermissionResults, err := s.rolePermissionRepository.GetUserRolePermission(ctx, userRoleIDs)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Login - Failed to get user role permission")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	userRoleMap := utils.MapUserRoleResponse(userRolePermissionResults)

	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Login - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::Login - Failed to rollback transaction")
			}
		}
	}()

	// Update last login at
	err = s.userRepository.UpdateUserLastLoginAt(ctx, tx, userResult.ID.String())
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Login - Failed to update last login at")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// generate UUID V7
	userFcmTokenID, err := utils.GenerateUUIDv7String()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Login - Failed to generate UUID V7")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// insert user fcm token
	err = s.userFcmTokenRepository.InsertNewUserFCMToken(ctx, tx, userFcmToken.UserFCMToken{
		ID:         userFcmTokenID,
		UserID:     userResult.ID,
		DeviceID:   req.DeviceID,
		DeviceType: strings.ToUpper(req.DeviceType),
		FcmToken:   req.FcmToken,
	})
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Login - Failed to insert new user fcm token")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// mapping login response data
	res = &dto.AuthUserResponse{
		Email: req.Email,
		EmailConfirmed: dto.EmailConfirmed{
			IsConfirm: true,
			CreatedAt: utils.FormatToWIB(*userResult.EmailVerifiedAt),
		},
		FullName:    userResult.FullName,
		PhoneNumber: *userProfileResult.PhoneNumber,
		Token: dto.TokenDetail{
			Token:     result.AccessToken,
			ExpiredAt: utils.FormatToWIB(result.TokenExpiredAt),
			CreatedAt: utils.FormatToWIB(result.CreatedAt),
			RefreshToken: dto.RefreshTokenDetail{
				RefreshToken: result.RefreshToken,
				ExpiredAt:    utils.FormatToWIB(result.RefreshTokenExpiredAt),
				CreatedAt:    utils.FormatToWIB(result.CreatedAt),
			},
		},
		UserRole: userRoleMap,
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Login - Failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// return response
	return res, nil
}

func (s *userService) Logout(ctx context.Context, accessToken string, locals *middleware.Locals) error {
	// parse token
	_, err := s.jwt.ParseTokenString(ctx, accessToken, locals.Username, locals.SessionID, constants.AccessTokenType)
	if err != nil {
		log.Error().Err(err).Any("access_token", accessToken).Msg("service::Logout - Failed to parse access token")
		return err_msg.NewCustomErrors(fiber.StatusUnauthorized, err_msg.WithMessage(constants.ErrInvalidAccessToken))
	}

	// delete user fcm token
	err = s.userFcmTokenRepository.DeleteUserFCMToken(ctx, locals.UserID)
	if err != nil {
		log.Error().Err(err).Any("access_token", accessToken).Msg("service::Logout - Failed to delete user fcm token")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// delete access token from redis
	userKey := fmt.Sprintf("%s:%s:%s", locals.Username, locals.SessionID, constants.AccessTokenType)
	err = s.redisRepository.Del(ctx, userKey)
	if err != nil {
		log.Error().Err(err).Any("access_token", accessToken).Msg("service::Logout - Failed to delete access token from Redis")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// delete refresh token from redis
	refreshTokenKey := fmt.Sprintf("%s:%s:%s", locals.Username, locals.SessionID, constants.RefreshTokenType)
	err = s.redisRepository.Del(ctx, refreshTokenKey)
	if err != nil {
		log.Error().Err(err).Any("access_token", accessToken).Msg("service::Logout - Failed to delete refresh token from Redis")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}

func (s *userService) GetCurrentUser(ctx context.Context, locals *middleware.Locals) (*dto.GetCurrentUserResponse, error) {
	// define variable
	var (
		res = new(dto.GetCurrentUserResponse)
	)

	// find user by id
	userResult, err := s.userRepository.FindByID(ctx, locals.UserID)
	if err != nil {
		log.Error().Err(err).Any("user_id", locals.UserID).Msg("service::GetCurrentUser - Failed to find user by id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get user profile
	userProfileResult, err := s.userProfileRepository.FindByUserID(ctx, userResult.ID.String())
	if err != nil {
		log.Error().Err(err).Any("payload", userResult.ID.String()).Msg("service::GetCurrentUser - Failed to find user profile by user id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get user role ids
	userRoleIDs, err := s.userRoleRepository.FindByUserID(ctx, userResult.ID.String())
	if err != nil {
		log.Error().Err(err).Any("payload", userResult.ID.String()).Msg("service::GetCurrentUser - Failed to find user role by user id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get user role permission
	userRolePermissionResults, err := s.rolePermissionRepository.GetUserRolePermission(ctx, userRoleIDs)
	if err != nil {
		log.Error().Err(err).Any("payload", userRoleIDs).Msg("service::GetCurrentUser - Failed to get user role permission")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	userRoleMap := utils.MapUserRoleResponse(userRolePermissionResults)

	// mapping get current user response data
	res = &dto.GetCurrentUserResponse{
		Email: userResult.Email,
		EmailConfirmed: dto.EmailConfirmed{
			IsConfirm: true,
			CreatedAt: utils.FormatToWIB(*userResult.EmailVerifiedAt),
		},
		FullName:    userResult.FullName,
		PhoneNumber: *userProfileResult.PhoneNumber,
		UserRole:    userRoleMap,
	}

	// return response
	return res, nil
}

func (s *userService) RefreshToken(ctx context.Context, accessToken string, locals *middleware.Locals) (*dto.AuthUserResponse, error) {
	// define variable
	var (
		res = new(dto.AuthUserResponse)
	)

	// parse token
	claims, err := s.jwt.ParseTokenString(ctx, accessToken, locals.Username, locals.SessionID, constants.RefreshTokenType)
	if err != nil {
		log.Error().Err(err).Any("access_token", accessToken).Msg("service::RefreshToken - Failed to parse access token")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnauthorized, err_msg.WithMessage(constants.ErrInvalidAccessToken))
	}

	// generate new token
	result, err := s.jwt.GenerateTokenString(ctx, jwt_handler.CostumClaimsPayload{
		UserID:     claims.UserID,
		Username:   claims.Username,
		Email:      claims.Email,
		FullName:   claims.FullName,
		SessionID:  utils.GenerateSessionUUID(),
		DeviceID:   claims.DeviceID,
		DeviceType: claims.DeviceType,
		FcmToken:   claims.FcmToken,
	})
	if err != nil {
		log.Error().Err(err).Any("payload", claims).Msg("service::RefreshToken - Failed to generate token string")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// find user by id
	userResult, err := s.userRepository.FindByID(ctx, claims.UserID)
	if err != nil {
		log.Error().Err(err).Any("user_id", claims.UserID).Msg("service::GetCurrentUser - Failed to find user by id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get user profile
	userProfileResult, err := s.userProfileRepository.FindByUserID(ctx, userResult.ID.String())
	if err != nil {
		log.Error().Err(err).Any("payload", userResult.ID.String()).Msg("service::RefreshToken - Failed to find user profile by user id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get user role ids
	userRoleIDs, err := s.userRoleRepository.FindByUserID(ctx, userResult.ID.String())
	if err != nil {
		log.Error().Err(err).Any("payload", userResult.ID.String()).Msg("service::RefreshToken - Failed to find user role by user id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get user role permission
	userRolePermissionResults, err := s.rolePermissionRepository.GetUserRolePermission(ctx, userRoleIDs)
	if err != nil {
		log.Error().Err(err).Any("payload", userRoleIDs).Msg("service::RefreshToken - Failed to get user role permission")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	userRoleMap := utils.MapUserRoleResponse(userRolePermissionResults)

	// mapping refresh token response data
	res = &dto.AuthUserResponse{
		Email: claims.Email,
		EmailConfirmed: dto.EmailConfirmed{
			IsConfirm: true,
			CreatedAt: utils.FormatToWIB(*userResult.EmailVerifiedAt),
		},
		FullName:    userResult.FullName,
		PhoneNumber: *userProfileResult.PhoneNumber,
		Token: dto.TokenDetail{
			Token:     result.AccessToken,
			ExpiredAt: utils.FormatToWIB(result.TokenExpiredAt),
			CreatedAt: utils.FormatToWIB(result.CreatedAt),
			RefreshToken: dto.RefreshTokenDetail{
				RefreshToken: result.RefreshToken,
				ExpiredAt:    utils.FormatToWIB(result.RefreshTokenExpiredAt),
				CreatedAt:    utils.FormatToWIB(result.CreatedAt),
			},
		},
		UserRole: userRoleMap,
	}

	// return response
	return res, nil
}

func (s *userService) ForgotPassword(ctx context.Context, req *dto.SendOtpNumberVerificationRequest) (*dto.ForgotPasswordResponse, error) {
	// declare variable
	var (
		res = new(dto.ForgotPasswordResponse)
	)

	// find user by email
	userResult, err := s.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::ForgotPassword - Failed to find user by email")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage("Failed to find user by email"))
	}

	// check user
	if userResult == nil {
		log.Error().Any("payload", req).Msg("service::ForgotPassword - User not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnprocessableEntity, err_msg.WithMessage(constants.ErrEmailNotRegistered))
	} else if userResult.EmailVerifiedAt == nil && userResult.OtpNumberVerifiedAt == nil {
		log.Error().Any("payload", req).Msg("service::ForgotPassword - User already verified")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnprocessableEntity, err_msg.WithMessage(constants.ErrEmailNotRegistered))
	}

	key := fmt.Sprintf("%s:%s", constants.ForgotPasswordKey, userResult.Email)

	_, err = s.redisRepository.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			generateSessionID := utils.GenerateSessionUUID()
			urlResetLink := fmt.Sprintf("%s/reset-password/%s", config.Envs.App.Domain, url.PathEscape(generateSessionID))

			err := s.redisRepository.Set(ctx, key, generateSessionID, 2*time.Minute)
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("service::ForgotPassword - Failed to set data redis")
				return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}

			// Send notification in background
			go func() {
				err := s.externalNotification.SendNotification(context.Background(), req.Email, constants.ForgotPasswordTemplateName, map[string]string{
					"full_name":      userResult.FullName,
					"url_reset_link": urlResetLink,
				})
				if err != nil {
					log.Error().Err(err).Msg("service::ForgotPassword - Failed to send notification in background")
				}
			}()

			res.Email = req.Email
			res.Sessions = generateSessionID

			return res, nil
		}

		log.Error().Err(err).Any("payload", req).Msg("service::ForgotPassword - Failed to get data redis")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	log.Warn().Any("payload", req).Msg("service::ForgotPassword - Too many requests for the same email, please wait until the current OTP expires")
	return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrTooManyReuqestOTPNumber))
}
