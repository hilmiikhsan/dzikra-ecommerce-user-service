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
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/province/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/rs/zerolog/log"
)

func (s *provinceService) GetListProvince(ctx context.Context) ([]dto.GetListProvinceResponse, error) {
	// 1) Try load from cache
	if cached, err := s.redisRepository.Get(ctx, constants.CacheKeyProvinces); err == nil && cached != "" {
		var out []dto.GetListProvinceResponse
		if err := json.Unmarshal([]byte(cached), &out); err == nil {
			log.Debug().Msg("service::GetListProvince - cache HIT")
			return out, nil
		}

		log.Warn().Err(err).Msg("service::GetListProvince - invalid cache, fetching fresh")
	}

	// cache not found or invalid, fetch from RajaOngkir API
	url := fmt.Sprintf("%s/province", config.Envs.RajaOngkir.BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListProvince - NewRequest failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// set headers
	req.Header.Set("key", config.Envs.RajaOngkir.ApiKey)

	// Send request
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListProvince - HTTP request failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer resp.Body.Close()

	// decode json response
	var payload rajaongkir.RajaOngkirProvincePayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("service::GetListProvince - Decode JSON failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// mapping response data
	out := make([]dto.GetListProvinceResponse, len(payload.Rajaongkir.Results))
	for i, data := range payload.Rajaongkir.Results {
		out[i] = dto.GetListProvinceResponse{
			ID:       data.ProvinceID,
			Province: data.Province,
		}
	}

	// set cache
	if bytes, err := json.Marshal(out); err == nil {
		if err := s.redisRepository.Set(ctx, constants.CacheKeyProvinces, string(bytes), constants.CacheTTL); err != nil {
			log.Warn().Err(err).Msg("service::GetListProvince - failed to set cache")
		}
	} else {
		log.Warn().Err(err).Msg("provinceService::GetListProvince - failed to marshal cache data")
	}

	return out, nil
}
