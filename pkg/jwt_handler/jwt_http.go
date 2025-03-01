package jwt_handler

import (
	"context"
	"fmt"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/infrastructure/config"
	redisPorts "github.com/Digitalkeun-Creative/be-dzikra-user-service/internal/infrastructure/redis/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-user-service/pkg/err_msg"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

var _ JWT = &jwtHandler{}

type jwtHandler struct {
	db redisPorts.RedisRepository
}

func NewJWT(db redisPorts.RedisRepository) *jwtHandler {
	return &jwtHandler{
		db: db,
	}
}

func (j *jwtHandler) GenerateTokenString(ctx context.Context, payload CostumClaimsPayload) (*GenerateTokenResponse, error) {
	// Hapus token lama berdasarkan pola lama (jika masih diperlukan, namun pada multiple session sebaiknya dihapus berdasarkan sessionID tertentu)
	for _, tokenType := range []string{constants.AccessTokenType, constants.RefreshTokenType} {
		key := fmt.Sprintf("%s:%s", payload.Username, tokenType)
		err := j.db.Del(ctx, key)
		if err != nil {
			log.Error().Err(err).Msgf("jwthandler::GenerateTokenString - Error while deleting %s from Redis", tokenType)
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	now := time.Now().UTC()

	// Parse access token duration from config
	accessTokenDuration, err := time.ParseDuration(config.Envs.Guard.JwtTokenExpiration)
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::GenerateTokenString - Error while parsing access token duration")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Parse refresh token duration from config
	refreshTokenDuration, err := time.ParseDuration(config.Envs.Guard.JwtRefreshTokenExpiration)
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::GenerateTokenString - Error while parsing refresh token duration")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Helper function untuk membuat token dengan claims, termasuk sessionID
	createToken := func(tokenType string, duration time.Duration) (string, time.Time, error) {
		expireTime := now.Add(duration)
		claims := CustomClaims{
			UserID:     payload.UserID,
			Username:   payload.Username,
			Email:      payload.Email,
			FullName:   payload.FullName,
			SessionID:  payload.SessionID,
			DeviceID:   payload.DeviceID,
			DeviceType: payload.DeviceType,
			FcmToken:   payload.FcmToken,
			CreatedAt:  now,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "user",
				Issuer:    config.Envs.App.Name,
				ExpiresAt: jwt.NewNumericDate(expireTime),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
		tokenString, err := token.SignedString([]byte(config.Envs.Guard.JwtPrivateKey))
		if err != nil {
			log.Error().Err(err).Msgf("jwthandler::GenerateTokenString - Error while signing %s", tokenType)
			return "", time.Time{}, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		return tokenString, expireTime, nil
	}

	// Generate access token use sessionID
	accessToken, accessExpireTime, err := createToken(constants.AccessTokenType, accessTokenDuration)
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::GenerateTokenString - Error while generating access token")
		return nil, err
	}

	// Save access token to Redis with unique key based on username and sessionID
	accessTokenKey := fmt.Sprintf("%s:%s:%s", payload.Username, payload.SessionID, constants.AccessTokenType)
	err = j.db.Set(ctx, accessTokenKey, accessToken, time.Until(accessExpireTime))
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::GenerateTokenString - Error while saving access token to Redis")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Generate refresh token use sessionID
	refreshToken, refreshExpireTime, err := createToken(constants.RefreshTokenType, refreshTokenDuration)
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::GenerateTokenString - Error while generating refresh token")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Save refresh token to Redis with unique key based on sessionID
	refreshTokenKey := fmt.Sprintf("%s:%s:%s", payload.Username, payload.SessionID, constants.RefreshTokenType)
	err = j.db.Set(ctx, refreshTokenKey, refreshToken, time.Until(refreshExpireTime))
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::GenerateTokenString - Error while saving refresh token to Redis")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Return the response
	return &GenerateTokenResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		TokenExpiredAt:        accessExpireTime,
		RefreshTokenExpiredAt: refreshExpireTime,
		CreatedAt:             now,
		SessionID:             payload.SessionID,
	}, nil
}

func (j *jwtHandler) ParseTokenString(ctx context.Context, tokenString, username, sessionID, tokenType string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	key := fmt.Sprintf("%s:%s:%s", username, sessionID, tokenType)

	storedToken, err := j.db.Get(ctx, key)
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::ParseTokenString - Token not found in Redis")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnauthorized, err_msg.WithMessage(constants.ErrTokenAlreadyExpired))
	}

	if storedToken != tokenString {
		log.Error().Msg("jwthandler::ParseTokenString - Token mismatch")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnauthorized, err_msg.WithMessage(constants.ErrTokenAlreadyExpired))
	}

	// Revoke the old refresh token after validation
	err = j.db.Del(ctx, key)
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::ParseTokenString - Error while deleting token from Redis")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Envs.Guard.JwtPrivateKey), nil
	})
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::ParseTokenString - Error while parsing token")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if !token.Valid {
		log.Error().Msg("jwthandler::ParseTokenString - Invalid token")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnauthorized, err_msg.WithMessage(constants.ErrTokenAlreadyExpired))
	}

	return claims, nil
}

func (j *jwtHandler) ParseMiddlewareTokenString(ctx context.Context, tokenString string) (*CustomClaims, error) {
	// define custom claims
	claims := &CustomClaims{}

	// Parse the JWT string and store the result in `claims`
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Envs.Guard.JwtPrivateKey), nil
	})
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::ParseMiddlewareTokenString - Error while parsing token")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnauthorized, err_msg.WithMessage(constants.ErrTokenAlreadyExpired))
	}

	// Check if the token is valid
	if !token.Valid {
		log.Error().Msg("jwthandler::ParseMiddlewareTokenString - Invalid token")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnauthorized, err_msg.WithMessage(constants.ErrTokenAlreadyExpired))
	}

	// Check if the token is stored in Redis
	key := fmt.Sprintf("%s:%s:%s", claims.Username, claims.SessionID, constants.AccessTokenType)
	storedToken, err := j.db.Get(ctx, key)
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::ParseMiddlewareTokenString - Token not found in Redis")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnauthorized, err_msg.WithMessage(constants.ErrTokenAlreadyExpired))
	}

	// Check if the token is match
	if storedToken != tokenString {
		log.Error().Msg("jwthandler::ParseMiddlewareTokenString - Token mismatch in Redis")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnauthorized, err_msg.WithMessage(constants.ErrTokenAlreadyExpired))
	}

	// If the token is valid, return the claims
	return claims, nil
}
