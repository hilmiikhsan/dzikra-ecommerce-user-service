package service

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (s *addressService) CreateAddress(ctx context.Context, req *dto.CreateOrUpdateAddressRequest) (*dto.CreateOrUpdateAddressResponse, error) {
	// get list province data
	provinceResults, err := s.provinceService.GetListProvince(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateAddress - failed to get province list")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check city id data
	isProvinceFound := false
	for _, prov := range provinceResults {
		if prov.ID == req.ProvinceID {
			isProvinceFound = true
			break
		}
	}
	if !isProvinceFound {
		log.Error().Str("province_id", req.ProvinceID).Msg("service::CreateAddress - invalid province_id")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("Invalid province_id"))
	}

	// convert provinceID to int
	provinceID, err := strconv.Atoi(req.ProvinceID)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateAddress - failed to convert province_id to int")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get list city data
	cityResults, err := s.cityService.GetListCity(ctx, provinceID)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateAddress - failed to get city list")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check city id data
	isCityFound := false
	for _, city := range cityResults {
		if city.ID == req.CityID {
			isCityFound = true
			break
		}
	}
	if !isCityFound {
		log.Error().Str("city_id", req.ProvinceID).Msg("service::CreateAddress - invalid city_id")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("Invalid city_id"))
	}

	// convert cityID to int
	cityID, err := strconv.Atoi(req.CityID)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateAddress - failed to convert city_id to int")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get list sub district data
	subDistrictResults, err := s.subDistrictService.GetListSubDistrict(ctx, cityID)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateAddress - failed to get sub district list")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check sub district id data
	isSubDistrictFound := false
	for _, subDistrict := range subDistrictResults {
		if subDistrict.ID == req.SubDistrictID {
			isSubDistrictFound = true
			break
		}
	}
	if !isSubDistrictFound {
		log.Error().Str("sub_district_id", req.SubDistrictID).Msg("service::CreateAddress - invalid sub_district_id")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("Invalid sub_district_id"))
	}

	// convert subDistrictID to int
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateAddress - failed to parse user_id")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateAddress - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::CreateAddress - Failed to rollback transaction")
			}
		}
	}()

	// create address
	res, err := s.addressRepository.InsertNewAddress(ctx, tx, &entity.Address{
		ReceivedName:        req.ReceivedName,
		Province:            req.Province,
		ProvinceVendorID:    req.ProvinceID,
		City:                req.City,
		CityVendorID:        req.CityID,
		SubDistrict:         req.SubDistrict,
		SubDistrictVendorID: req.SubDistrictID,
		Address:             req.Address,
		PostalCode:          req.PostalCode,
		UserID:              userID,
	})
	if err != nil {
		log.Error().Err(err).Msg("service::CreateAddress - failed to create address")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::CreateAddress - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrUpdateAddressResponse{
		ID:                  res.ID,
		Province:            res.Province,
		ProvinceVendorID:    res.ProvinceVendorID,
		City:                res.City,
		CityVendorID:        res.CityVendorID,
		District:            utils.NullStringPtr(res.District),
		SubDistrict:         res.SubDistrict,
		SubDistrictVendorID: res.SubDistrictVendorID,
		Address:             res.Address,
		PostalCode:          res.PostalCode,
		ReceivedName:        res.ReceivedName,
		UserID:              req.UserID,
	}, nil
}

func (s addressService) UpdateAddress(ctx context.Context, req *dto.CreateOrUpdateAddressRequest, addressID int) (*dto.CreateOrUpdateAddressResponse, error) {
	// get list province data
	provinceResults, err := s.provinceService.GetListProvince(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateAddress - failed to get province list")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check city id data
	isProvinceFound := false
	for _, prov := range provinceResults {
		if prov.ID == req.ProvinceID {
			isProvinceFound = true
			break
		}
	}
	if !isProvinceFound {
		log.Error().Str("province_id", req.ProvinceID).Msg("service::UpdateAddress - invalid province_id")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("Invalid province_id"))
	}

	// convert provinceID to int
	provinceID, err := strconv.Atoi(req.ProvinceID)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateAddress - failed to convert province_id to int")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get list city data
	cityResults, err := s.cityService.GetListCity(ctx, provinceID)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateAddress - failed to get city list")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check city id data
	isCityFound := false
	for _, city := range cityResults {
		if city.ID == req.CityID {
			isCityFound = true
			break
		}
	}
	if !isCityFound {
		log.Error().Str("city_id", req.ProvinceID).Msg("service::UpdateAddress - invalid city_id")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("Invalid city_id"))
	}

	// convert cityID to int
	cityID, err := strconv.Atoi(req.CityID)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateAddress - failed to convert city_id to int")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get list sub district data
	subDistrictResults, err := s.subDistrictService.GetListSubDistrict(ctx, cityID)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateAddress - failed to get sub district list")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check sub district id data
	isSubDistrictFound := false
	for _, subDistrict := range subDistrictResults {
		if subDistrict.ID == req.SubDistrictID {
			isSubDistrictFound = true
			break
		}
	}
	if !isSubDistrictFound {
		log.Error().Str("sub_district_id", req.SubDistrictID).Msg("service::UpdateAddress - invalid sub_district_id")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("Invalid sub_district_id"))
	}

	// convert subDistrictID to int
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateAddress - failed to parse user_id")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::UpdateAddress - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::UpdateAddress - Failed to rollback transaction")
			}
		}
	}()

	// update address
	res, err := s.addressRepository.UpdateAddress(ctx, tx, addressID, &entity.Address{
		ReceivedName:        req.ReceivedName,
		Province:            req.Province,
		ProvinceVendorID:    req.ProvinceID,
		City:                req.City,
		CityVendorID:        req.CityID,
		SubDistrict:         req.SubDistrict,
		SubDistrictVendorID: req.SubDistrictID,
		Address:             req.Address,
		PostalCode:          req.PostalCode,
		UserID:              userID,
	})
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrAddressNotFound) {
			log.Error().Err(err).Msg("service::UpdateAddress - address not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrAddressNotFound))
		}

		log.Error().Err(err).Msg("service::UpdateAddress - failed to update address")
		return nil, err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::UpdateAddress - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrUpdateAddressResponse{
		ID:                  res.ID,
		Province:            res.Province,
		ProvinceVendorID:    res.ProvinceVendorID,
		City:                res.City,
		CityVendorID:        res.CityVendorID,
		District:            utils.NullStringPtr(res.District),
		SubDistrict:         res.SubDistrict,
		SubDistrictVendorID: res.SubDistrictVendorID,
		Address:             res.Address,
		PostalCode:          res.PostalCode,
		ReceivedName:        res.ReceivedName,
		UserID:              req.UserID,
	}, nil
}

func (s *addressService) RemoveAddress(ctx context.Context, addressID int, userID string) error {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::RemoveAddress - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::RemoveAddress - Failed to rollback transaction")
			}
		}
	}()

	// convert userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Error().Err(err).Msg("service::RemoveAddress - failed to parse user_id")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// soft delete address
	err = s.addressRepository.SoftDeleteAddressByID(ctx, tx, addressID, userUUID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrAddressNotFound) {
			log.Error().Err(err).Msg("service::RemoveAddress - address not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrAddressNotFound))
		}

		log.Error().Err(err).Msg("service::RemoveAddress - failed to soft delete address")
		return err_msg.NewCustomErrors(http.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::RemoveAddress - failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}

func (s *addressService) GetListAddress(ctx context.Context, userID string) ([]dto.GetListAddressResponse, error) {
	// convert userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListAddress - failed to parse user_id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get list address
	addresses, err := s.addressRepository.FindAllAddressByUserID(ctx, userUUID)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListAddress - failed to get list address")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if addresses is nil
	if addresses == nil {
		addresses = []dto.GetListAddressResponse{}
	}

	return addresses, nil
}

func (s *addressService) GetDetailAddress(ctx context.Context, addressID int, userID string) (*dto.GetListAddressResponse, error) {
	// convert userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Error().Err(err).Msg("service::GetDetailAddress - failed to parse user_id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// get detail address
	address, err := s.addressRepository.FindDetailAddressByID(ctx, addressID, userUUID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrAddressNotFound) {
			log.Error().Err(err).Msg("service::GetDetailAddress - address not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrAddressNotFound))
		}

		log.Error().Err(err).Msg("service::GetDetailAddress - failed to get detail address")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.GetListAddressResponse{
		ID:                  address.ID,
		Province:            address.Province,
		City:                address.City,
		SubDistrict:         address.SubDistrict,
		Address:             address.Address,
		PostalCode:          address.PostalCode,
		CityVendorID:        address.CityVendorID,
		ProvinceVendorID:    address.ProvinceVendorID,
		SubDistrictVendorID: address.SubDistrictVendorID,
		ReceivedName:        address.ReceivedName,
	}, nil
}

func (s *addressService) GetAddressesByIds(ctx context.Context, ids []int64) ([]dto.GetAddressesByIdsResponse, error) {
	addresses, err := s.addressRepository.FindAddressesByIds(ctx, ids)
	if err != nil {
		log.Error().Err(err).Msg("service::GetAddressesByIds - error fetching addresses by IDs")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	var result []dto.GetAddressesByIdsResponse
	for _, address := range addresses {
		var district string
		if address.District.Valid {
			district = address.District.String
		}

		result = append(result, dto.GetAddressesByIdsResponse{
			ID:                  address.ID,
			Province:            address.Province,
			City:                address.City,
			SubDistrict:         address.SubDistrict,
			CityVendorID:        address.CityVendorID,
			ProvinceVendorID:    address.ProvinceVendorID,
			SubDistrictVendorID: address.SubDistrictVendorID,
			Address:             address.Address,
			PostalCode:          address.PostalCode,
			ReceivedName:        address.ReceivedName,
			District:            district,
			UserID:              address.UserID.String(),
		})
	}

	return result, nil
}
