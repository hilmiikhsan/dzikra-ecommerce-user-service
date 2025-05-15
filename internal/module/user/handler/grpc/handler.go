package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/cmd/proto/tokenvalidation"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/rs/zerolog/log"
)

func (h *userGrpcAPI) ValidateToken(ctx context.Context, req *tokenvalidation.TokenRequest) (*tokenvalidation.TokenResponse, error) {
	var (
		token = req.GetToken()
	)

	if token == "" {
		err := fmt.Errorf("grpcHandler::ValidateToken - Token is empty")
		log.Error().Err(nil).Msg(err.Error())

		return &tokenvalidation.TokenResponse{
			Message: err.Error(),
		}, nil
	}

	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimSpace(token[7:])
	}

	claims, err := h.jwt.ParseMiddlewareTokenString(ctx, token)
	if err != nil {
		log.Error().Err(err).Any("token", token).Msg("grpcHandler::ValidateToken - error parsing token")

		return &tokenvalidation.TokenResponse{
			Message: err.Error(),
		}, nil
	}

	var userRoles []*tokenvalidation.ApplicationPermission
	for _, appPerm := range claims.UserRoles {
		var perms []*tokenvalidation.UserRoleAppPermission
		for _, p := range appPerm.ApplicationPermission {
			perms = append(perms, &tokenvalidation.UserRoleAppPermission{
				ApplicationId: p.ApplicationID,
				Name:          p.Name,
				Permissions:   p.Permissions,
			})
		}

		userRoles = append(userRoles, &tokenvalidation.ApplicationPermission{
			ApplicationPermissions: perms,
			Roles:                  appPerm.Roles,
		})
	}

	return &tokenvalidation.TokenResponse{
		Message: constants.SuccessMessage,
		Data: &tokenvalidation.UserData{
			UserId:    claims.UserID,
			Email:     claims.Email,
			FullName:  claims.FullName,
			UserRoles: userRoles,
		},
	}, nil
}
