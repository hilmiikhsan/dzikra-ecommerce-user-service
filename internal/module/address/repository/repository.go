package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.AddressRepository = &addressRepository{}

type addressRepository struct {
	db *sqlx.DB
}

func NewAddressRepository(db *sqlx.DB) *addressRepository {
	return &addressRepository{
		db: db,
	}
}

func (r *addressRepository) InsertNewAddress(ctx context.Context, tx *sqlx.Tx, data *entity.Address) (*entity.Address, error) {
	var res = new(entity.Address)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryInsertNewAddress),
		data.Province,
		data.City,
		data.District,
		data.SubDistrict,
		data.PostalCode,
		data.Address,
		data.ReceivedName,
		data.UserID,
		data.CityVendorID,
		data.ProvinceVendorID,
		data.SubDistrictVendorID,
	).Scan(
		&res.ID,
		&res.Province,
		&res.City,
		&res.District,
		&res.SubDistrict,
		&res.PostalCode,
		&res.Address,
		&res.ReceivedName,
		&res.UserID,
		&res.CityVendorID,
		&res.ProvinceVendorID,
		&res.SubDistrictVendorID,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewAddress - Failed to insert new address")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}

func (r *addressRepository) UpdateAddress(ctx context.Context, tx *sqlx.Tx, id int, data *entity.Address) (*entity.Address, error) {
	var res = new(entity.Address)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryUpdateAddress),
		data.Province,
		data.City,
		data.District,
		data.SubDistrict,
		data.PostalCode,
		data.Address,
		data.ReceivedName,
		data.UserID,
		data.CityVendorID,
		data.ProvinceVendorID,
		data.SubDistrictVendorID,
		id,
	).Scan(
		&res.ID,
		&res.Province,
		&res.City,
		&res.District,
		&res.SubDistrict,
		&res.PostalCode,
		&res.Address,
		&res.ReceivedName,
		&res.UserID,
		&res.CityVendorID,
		&res.ProvinceVendorID,
		&res.SubDistrictVendorID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::UpdateAddress - address with id %d is not found", id)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrAddressNotFound)
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateAddress - Failed to update address")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}

func (r *addressRepository) SoftDeleteAddressByID(ctx context.Context, tx *sqlx.Tx, id int, userID uuid.UUID) error {
	result, err := tx.ExecContext(ctx, r.db.Rebind(querySoftDeleteAddressByID), id, userID)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("repository::SoftDeleteAddressByID - Failed to soft delete voucher")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteAddressByID - Failed to fetch rows affected")
		return err
	}

	if rowsAffected == 0 {
		errNotFound := errors.New(constants.ErrAddressNotFound)
		log.Error().Err(errNotFound).Int("id", id).Msg("repository::SoftDeleteAddressByID - Voucher not found")
		return errNotFound
	}

	return nil
}

func (r *addressRepository) FindAllAddressByUserID(ctx context.Context, userID uuid.UUID) ([]dto.GetListAddressResponse, error) {
	var responses []entity.Address

	if err := r.db.SelectContext(ctx, &responses, r.db.Rebind(queryFindAllAddressByUserID), userID); err != nil {
		log.Error().Err(err).Msg("repository::FindAllAddressByUserID - error executing query")
		return nil, err
	}

	addresses := make([]dto.GetListAddressResponse, 0, len(responses))
	for _, v := range responses {
		addresses = append(addresses, dto.GetListAddressResponse{
			ID:                  v.ID,
			Province:            v.Province,
			City:                v.City,
			SubDistrict:         v.SubDistrict,
			CityVendorID:        v.CityVendorID,
			ProvinceVendorID:    v.ProvinceVendorID,
			SubDistrictVendorID: v.SubDistrictVendorID,
			Address:             v.Address,
			PostalCode:          v.PostalCode,
			ReceivedName:        v.ReceivedName,
		})
	}

	return addresses, nil
}
