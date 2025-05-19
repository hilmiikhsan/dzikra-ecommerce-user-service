package service

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/shipping/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (s *shippingService) CalculateShippingCost(ctx context.Context, req *dto.CalculateShippingCostRequest, userID string) ([]dto.CalculateShippingCostResponse, error) {
	// convert user id to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Error().Err(err).Msg("service::CalculateShippingCost - Failed to parse user ID")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get list cart item
	cartItems, err := s.cartRepository.FindListCartByUserID(ctx, userUUID)
	if err != nil {
		log.Error().Err(err).Msg("service::CalculateShippingCost - Failed to get list cart item")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if cart item is empty
	if len(cartItems) == 0 {
		log.Error().Msg("service::CalculateShippingCost - Cart item is empty")
		return nil, err_msg.NewCustomErrors(http.StatusBadRequest, err_msg.WithMessage(constants.ErrCartItemIsEmpty))
	}

	// sum total weight (in kg) and convert to grams
	var totalWeightKg float64
	for _, item := range cartItems {
		var weightKg float64

		if item.ProductVariantID != 0 {
			weightKg = item.ProductVariantWeight
		} else {
			weightKg = item.ProductWeight
		}

		totalWeightKg += weightKg * float64(item.Quantity)
	}

	totalWeightGrams := strconv.Itoa(int(totalWeightKg * 1000))

	// Fetch the users address
	address, err := s.addressRepository.FindDetailAddressByID(ctx, req.AddressID, userUUID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrAddressNotFound) {
			log.Warn().Err(err).Msg("service::CalculateShippingCost - Address not found")
			return nil, err_msg.NewCustomErrors(http.StatusBadRequest, err_msg.WithMessage(constants.ErrAddressNotFound))
		}

		log.Error().Err(err).Msg("service::CalculateShippingCost - Failed to get address")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// call RajaOngkir integration
	costs, err := s.rajaongkirService.GetShippingCost(ctx, totalWeightGrams, req.Courier, address)
	if err != nil {
		log.Error().Err(err).Msg("service::CalculateShippingCost - Failed to get shipping cost")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if costs is empty
	if len(costs) == 0 || len(costs[0].Cost) == 0 {
		log.Error().Msg("service::CalculateShippingCost - No shipping options returned")
		return nil, err_msg.NewCustomErrors(http.StatusBadRequest, err_msg.WithMessage(constants.ErrNoShippingOptionReturned))
	}

	// map into module DTO
	res := make([]dto.CalculateShippingCostResponse, len(costs))
	for i, c := range costs {
		res[i] = dto.CalculateShippingCostResponse{
			Code:  c.Code,
			Name:  c.Name,
			Costs: make([]dto.DetailShippingCost, len(c.Cost)),
		}
		for j, svc := range c.Cost {
			ds := dto.DetailShippingCost{
				Service:     svc.Service,
				Description: svc.Description,
				Cost:        make([]dto.DetailCost, len(svc.Cost)),
			}
			for k, raw := range svc.Cost {
				ds.Cost[k] = dto.DetailCost{
					Value: raw.Value,
					Etd:   raw.Etd,
					Note:  raw.Note,
				}
			}
			res[i].Costs[j] = ds
		}
	}

	return res, nil
}
