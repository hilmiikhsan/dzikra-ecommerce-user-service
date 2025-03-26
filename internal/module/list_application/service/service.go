package service

import (
	"context"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/list_application/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/list_application/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *applicationService) GetListApplication(ctx context.Context) ([]dto.GetListApplicationResponse, error) {
	// get list application
	appEntities, err := s.applicationRepository.FindAllApplication(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListApplication - Failed to get list application")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// validate if appEntities is nil
	if appEntities == nil {
		appEntities = []entity.Application{}
	}

	// create mapping response
	var responses []dto.GetListApplicationResponse
	for _, app := range appEntities {
		responses = append(responses, dto.GetListApplicationResponse{
			ID:   app.ID,
			Name: app.Name,
		})
	}

	// return response
	return responses, nil
}

func (s *applicationService) GetListPermissionByApp(ctx context.Context, appIDsParam string) (*dto.GetListPermissionByAppResponse, error) {
	// declare variables
	var apps []dto.PermissionApp
	var err error

	// check if appIDsParam is not empty
	if strings.TrimSpace(appIDsParam) != "" {
		ids := strings.Split(appIDsParam, ",")
		for i, id := range ids {
			ids[i] = strings.TrimSpace(id)
		}

		// get permission apps by IDs
		apps, err = s.applicationRepository.FindPermissionAppsByIDs(ctx, ids)
		if err != nil {
			log.Error().Err(err).Msg("service::GetListPermissionByApp - Failed to get permission apps by IDs")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage("Internal server error"))
		}
	} else {
		apps = []dto.PermissionApp{}
	}

	// validate if apps is nil
	if apps == nil {
		apps = []dto.PermissionApp{}
	}

	// create mapping response
	response := &dto.GetListPermissionByAppResponse{
		ApplicationPermissions: apps,
	}

	// return response
	return response, nil
}
