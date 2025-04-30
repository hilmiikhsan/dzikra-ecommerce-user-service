package service

import (
	"context"
	"net/http"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/city/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/rs/zerolog/log"
)

func (s *cityService) GetListCity(ctx context.Context, provinceID int) ([]dto.GetListCityResponse, error) {
	// get lis city
	cityResults, err := s.rajaongkirService.GetListCity(ctx, provinceID)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListCity - Failed to get list city")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return cityResults, nil
}
