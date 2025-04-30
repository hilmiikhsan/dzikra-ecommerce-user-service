package service

import (
	"context"
	"net/http"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/sub_district/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/rs/zerolog/log"
)

func (s *subDistrictService) GetListSubDistrict(ctx context.Context, districtID int) ([]dto.GetListSubDistrictResponse, error) {
	// get list dub district
	subDistrictResults, err := s.rajaongkirService.GetListSubDistrict(ctx, districtID)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListSubDistrict - Failed to get list sub district")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return subDistrictResults, nil
}
