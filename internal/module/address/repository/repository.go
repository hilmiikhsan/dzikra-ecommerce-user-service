package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/gofiber/fiber/v2"
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
