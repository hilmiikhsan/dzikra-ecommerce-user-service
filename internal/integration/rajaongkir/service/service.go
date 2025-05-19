package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/integration/rajaongkir/dto"
	address "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/entity"
	city "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/city/dto"
	province "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/province/dto"
	subDistrict "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/sub_district/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/rs/zerolog/log"
)

func (s *rajaongkirService) GetListCity(ctx context.Context, provinceID int) ([]city.GetListCityResponse, error) {
	// declare key cache list city
	key := fmt.Sprintf("%s:%d", constants.CacheKeyCitys, provinceID)

	// Try load from cache
	if cached, err := s.redisRepository.Get(ctx, key); err == nil && cached != "" {
		var out []city.GetListCityResponse
		if err := json.Unmarshal([]byte(cached), &out); err == nil {
			log.Debug().Msg("service::GetListCity - cache HIT")
			return out, nil
		}

		log.Warn().Err(err).Msg("service::GetListCity - invalid cache, fetching fresh")
	}

	// cache not found or invalid, fetch from RajaOngkir API
	url := fmt.Sprintf("%s/city?province=%d", config.Envs.RajaOngkir.BaseURL, provinceID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListCity - NewRequest failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// set headers
	req.Header.Set("key", config.Envs.RajaOngkir.ApiKey)

	// Send request
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListCity - HTTP request failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer resp.Body.Close()

	// decode json response
	var payload dto.RajaOngkirCityPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("service::GetListCity - Decode JSON failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// mapping response data
	out := make([]city.GetListCityResponse, len(payload.Rajaongkir.Results))
	for i, data := range payload.Rajaongkir.Results {
		out[i] = city.GetListCityResponse{
			ID:           data.CityID,
			City:         data.CityName,
			Type:         data.Type,
			ProvinceName: data.ProvinceName,
			PostalCode:   data.PostalCode,
		}
	}

	// set cache
	if bytes, err := json.Marshal(out); err == nil {
		if err := s.redisRepository.Set(ctx, key, string(bytes), constants.CacheTTL); err != nil {
			log.Warn().Err(err).Msg("service::GetListCity - failed to set cache")
		}
	} else {
		log.Warn().Err(err).Msg("service::GetListCity - failed to marshal cache data")
	}

	return out, nil
}

func (s *rajaongkirService) GetListProvince(ctx context.Context) ([]province.GetListProvinceResponse, error) {
	// Try load from cache
	if cached, err := s.redisRepository.Get(ctx, constants.CacheKeyProvinces); err == nil && cached != "" {
		var out []province.GetListProvinceResponse
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
	var payload dto.RajaOngkirProvincePayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("service::GetListProvince - Decode JSON failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// mapping response data
	out := make([]province.GetListProvinceResponse, len(payload.Rajaongkir.Results))
	for i, data := range payload.Rajaongkir.Results {
		out[i] = province.GetListProvinceResponse{
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

func (s *rajaongkirService) GetListSubDistrict(ctx context.Context, districtID int) ([]subDistrict.GetListSubDistrictResponse, error) {
	// declare key cache list sub district
	key := fmt.Sprintf("%s:%d", constants.CacheKeySubDistricts, districtID)

	// Try load from cache
	if cached, err := s.redisRepository.Get(ctx, key); err == nil && cached != "" {
		var out []subDistrict.GetListSubDistrictResponse
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
	var payload dto.RajaOngkirSubDistrictPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("service::GetListSubDistrict - Decode JSON failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// mapping response data
	out := make([]subDistrict.GetListSubDistrictResponse, len(payload.Rajaongkir.Results))
	for i, data := range payload.Rajaongkir.Results {
		out[i] = subDistrict.GetListSubDistrictResponse{
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

func (s *rajaongkirService) GetShippingCost(ctx context.Context, weight, courier string, address *address.Address) ([]dto.CostResult, error) {
	key := fmt.Sprintf("%s:%d:%s:%s",
		constants.CacheKeyShippingCost,
		address.ID,
		courier,
		weight,
	)

	if cached, err := s.redisRepository.Get(ctx, key); err == nil && cached != "" {
		var out []dto.CostResult

		if err := json.Unmarshal([]byte(cached), &out); err == nil {
			log.Debug().Msg("service::GetShippingCost - cache HIT")
			return out, nil
		}

		log.Warn().Err(err).Msg("service::GetShippingCost - invalid cache, fetching fresh")
	}

	values := url.Values{}

	if config.Envs.RajaOngkir.ApiKeyType == "BASIC" {
		values.Set("origin", config.Envs.RajaOngkir.OriginCityID)
		values.Set("destination", address.CityVendorID)
	} else {
		values.Set("origin", config.Envs.RajaOngkir.OriginSubDistrictID)
		values.Set("originType", "subdistrict")
		values.Set("destination", address.SubDistrictVendorID)
		values.Set("destinationType", "subdistrict")
	}

	values.Set("weight", weight)
	values.Set("courier", courier)

	endpoint := fmt.Sprintf("%s/cost", config.Envs.RajaOngkir.BaseURL)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBufferString(values.Encode()))
	if err != nil {
		log.Error().Err(err).Msg("service::CalculateShippingCost - failed to build request")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("key", config.Envs.RajaOngkir.ApiKey)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Error().Err(err).Msg("service::CalculateShippingCost - API call failed")
		return nil, err_msg.NewCustomErrors(http.StatusBadGateway, err_msg.WithMessage(constants.ErrExternalServiceUnavailable))
	}
	defer resp.Body.Close()

	var payload dto.RajaOngkirCostPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("rajaongkir::GetShippingCost - decode failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	results := payload.Rajaongkir.Results

	if data, err := json.Marshal(results); err == nil {
		if err := s.redisRepository.Set(ctx, key, string(data), constants.CacheShippingCostTTL); err != nil {
			log.Warn().Err(err).Msg("service::GetShippingCost - failed to set cache")
		}
	} else {
		log.Warn().Err(err).Msg("service::GetShippingCost - failed to marshal cache data")
	}

	return results, nil
}

func (s *rajaongkirService) GetWaybill(ctx context.Context, waybill, courier string) (*dto.GetWaybillResponse, error) {
	values := url.Values{}
	values.Set("waybill", waybill)
	values.Set("courier", courier)

	endpoint := fmt.Sprintf("%s/waybill", config.Envs.RajaOngkir.BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBufferString(values.Encode()))
	if err != nil {
		log.Error().Err(err).Msg("service::GetWaybill - NewRequest failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("key", config.Envs.RajaOngkir.ApiKey)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("service::GetWaybill - HTTP request failed")
		return nil, err_msg.NewCustomErrors(http.StatusBadGateway, err_msg.WithMessage(constants.ErrExternalServiceUnavailable))
	}
	defer resp.Body.Close()

	var payload dto.RajaOngkirWaybillPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("service::GetWaybill - JSON decode failed")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	r := payload.Rajaongkir.Result

	out := &dto.GetWaybillResponse{
		Summary: dto.WaybillSummary{
			Resi:         r.Summary.WaybillNumber,
			ServiceCode:  r.Summary.ServiceCode,
			WaybillDate:  r.Summary.WaybillDate,
			ShipperName:  r.Summary.ShipperName,
			ReceiverName: r.Summary.ReceiverName,
			Origin:       r.Summary.Origin,
			Destination:  r.Summary.Destination,
			Status:       r.Summary.Status,
			CourierName:  r.Summary.CourierName,
		},
		Manifest: make([]dto.WaybillManifest, len(r.Manifest)),
		DeliveryStatus: dto.WaybillDeliveryStatus{
			Status:      r.DeliveryStatus.Status,
			PODReceiver: r.DeliveryStatus.PODReceiver,
			PODDate:     r.DeliveryStatus.PODDate,
			PODTime:     r.DeliveryStatus.PODTime,
		},
	}

	for i, m := range r.Manifest {
		out.Manifest[i] = dto.WaybillManifest{
			Description: m.ManifestDescription,
			Date:        m.ManifestDate,
			Time:        m.ManifestTime,
			CityName:    m.CityName,
		}
	}

	return out, nil
}
