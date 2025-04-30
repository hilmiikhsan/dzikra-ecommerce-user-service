package service

import (
	"context"
	"net/http"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/province/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/rs/zerolog/log"
)

func (s *provinceService) GetListProvince(ctx context.Context) ([]dto.GetListProvinceResponse, error) {
	// get list province
	provinceResults, err := s.rajaongkirService.GetListProvince(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListProvince - Failed to get list province")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return provinceResults, nil
}
