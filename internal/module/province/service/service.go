package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/province/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/rs/zerolog/log"
)

const (
	cacheKeyProvinces = "address:provinces"
	cacheTTL          = 10 * time.Minute
)

func (s *provinceService) GetListProvince(ctx context.Context) ([]dto.GetListProvinceResponse, error) {
	// 1) Try load from cache
	if cached, err := s.redisRepository.Get(ctx, cacheKeyProvinces); err == nil && cached != "" {
		var out []dto.GetListProvinceResponse
		if err := json.Unmarshal([]byte(cached), &out); err == nil {
			log.Debug().Msg("provinceService::GetListProvince - cache HIT")
			return out, nil
		}

		log.Warn().Err(err).Msg("provinceService::GetListProvince - invalid cache, fetching fresh")
	}

	// cache not found or invalid, fetch from RajaOngkir API
	url := fmt.Sprintf("%s/province", config.Envs.RajaOngkir.BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("provinceService::GetListProvince - NewRequest failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// set headers
	req.Header.Set("key", config.Envs.RajaOngkir.ApiKey)

	// Send request
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("provinceService::GetListProvince - HTTP request failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer resp.Body.Close()

	// decode json response
	var payload struct {
		Rajaongkir struct {
			Results []struct {
				ProvinceID string `json:"province_id"`
				Province   string `json:"province"`
			} `json:"results"`
		} `json:"rajaongkir"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("provinceService::GetListProvince - Decode JSON failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// mapping response data
	out := make([]dto.GetListProvinceResponse, len(payload.Rajaongkir.Results))
	for i, r := range payload.Rajaongkir.Results {
		out[i] = dto.GetListProvinceResponse{
			ID:       r.ProvinceID,
			Province: r.Province,
		}
	}

	// set cache
	if bytes, err := json.Marshal(out); err == nil {
		if err := s.redisRepository.Set(ctx, cacheKeyProvinces, string(bytes), cacheTTL); err != nil {
			log.Warn().Err(err).Msg("provinceService::GetListProvince - failed to set cache")
		}
	} else {
		log.Warn().Err(err).Msg("provinceService::GetListProvince - failed to marshal cache data")
	}

	return out, nil
}
