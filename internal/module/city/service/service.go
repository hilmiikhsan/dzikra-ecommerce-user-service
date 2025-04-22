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
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/city/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/rs/zerolog/log"
)

func (s *cityService) GetListCity(ctx context.Context, provinceID int) ([]dto.GetCityResponse, error) {
	// 1) Try load from cache
	if cached, err := s.redisRepository.Get(ctx, constants.CacheKeyCitys); err == nil && cached != "" {
		var out []dto.GetCityResponse
		if err := json.Unmarshal([]byte(cached), &out); err == nil {
			log.Debug().Msg("cityService::GetCityResponse - cache HIT")
			return out, nil
		}

		log.Warn().Err(err).Msg("cityService::GetCityResponse - invalid cache, fetching fresh")
	}

	// cache not found or invalid, fetch from RajaOngkir API
	url := fmt.Sprintf("%s/city?province=%d", config.Envs.RajaOngkir.BaseURL, provinceID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("cityService::GetCityResponse - NewRequest failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// set headers
	req.Header.Set("key", config.Envs.RajaOngkir.ApiKey)

	// Send request
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("cityService::GetCityResponse - HTTP request failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer resp.Body.Close()

	// decode json response
	var payload rajaongkir.RajaOngkirCityPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("cityService::GetCityResponse - Decode JSON failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// mapping response data
	out := make([]dto.GetCityResponse, len(payload.Rajaongkir.Results))
	for i, data := range payload.Rajaongkir.Results {
		out[i] = dto.GetCityResponse{
			ID:           data.CityID,
			City:         data.CityName,
			Type:         data.Type,
			ProvinceName: data.ProvinceName,
			PostalCode:   data.PostalCode,
		}
	}

	// set cache
	if bytes, err := json.Marshal(out); err == nil {
		if err := s.redisRepository.Set(ctx, constants.CacheKeyCitys, string(bytes), constants.CacheTTL); err != nil {
			log.Warn().Err(err).Msg("cityService::GetCityResponse - failed to set cache")
		}
	} else {
		log.Warn().Err(err).Msg("cityService::GetCityResponse - failed to marshal cache data")
	}

	return out, nil
}
