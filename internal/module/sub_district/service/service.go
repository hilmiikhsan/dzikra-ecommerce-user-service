package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/config"
	rajaongkir "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/rajaongkir/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/sub_district/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/rs/zerolog/log"
)

func (s *subDistrictService) GetListSubDistrict(ctx context.Context, districtID int) ([]dto.GetListSubDistrictResponse, error) {
	// declare key cache list sub district
	key := fmt.Sprintf("%s:%d", constants.CacheKeySubDistricts, districtID)

	// 1) Try load from cache
	if cached, err := s.redisRepository.Get(ctx, key); err == nil && cached != "" {
		var out []dto.GetListSubDistrictResponse
		if err := json.Unmarshal([]byte(cached), &out); err == nil {
			log.Debug().Msg("service::GetListSubDistrict - cache HIT")
			return out, nil
		}

		log.Warn().Err(err).Msg("service::GetListSubDistrict - invalid cache, fetching fresh")
	}

	// cache not found or invalid, fetch from RajaOngkir API
	url := fmt.Sprintf("%s/subdistrict?city=%d", config.Envs.RajaOngkir.BaseURL, districtID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListSubDistrict - NewRequest failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// set headers
	req.Header.Set("key", config.Envs.RajaOngkir.ApiKey)

	// Send request
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListSubDistrict - HTTP request failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer resp.Body.Close()

	// decode json response
	var payload rajaongkir.RajaOngkirSubDistrictPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("service::GetListSubDistrict - Decode JSON failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// mapping response data
	out := make([]dto.GetListSubDistrictResponse, len(payload.Rajaongkir.Results))
	for i, data := range payload.Rajaongkir.Results {
		out[i] = dto.GetListSubDistrictResponse{
			ID:              data.SubDistrictID,
			CityName:        data.CityName,
			CityType:        data.CityType,
			ProvinceName:    data.ProvinceName,
			SubDistrictName: data.SubDistrictName,
		}
	}

	// set cache
	if bytes, err := json.Marshal(out); err == nil {
		if err := s.redisRepository.Set(ctx, key, string(bytes), constants.CacheTTL); err != nil {
			log.Warn().Err(err).Msg("service::GetListSubDistrict - failed to set cache")
		}
	} else {
		log.Warn().Err(err).Msg("service::GetListSubDistrict - failed to marshal cache data")
	}

	return out, nil
}
