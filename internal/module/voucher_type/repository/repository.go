package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher_type/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher_type/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.VoucherTypeRepository = &voucherTypeRepository{}

type voucherTypeRepository struct {
	db *sqlx.DB
}

func NewVoucherTypeRepository(db *sqlx.DB) *voucherTypeRepository {
	return &voucherTypeRepository{
		db: db,
	}
}

func (r *voucherTypeRepository) CountVoucherTypeByType(ctx context.Context, voucherType string) (*entity.VoucherType, error) {
	var res entity.VoucherType

	err := r.db.GetContext(ctx, &res, r.db.Rebind(queryCountVoucherTypeByType), voucherType)
	if err != nil {
		if strings.Contains(err.Error(), "invalid input value for enum voucher_type_enum") {
			log.Error().Err(err).Str("voucherType", voucherType).Msg("repository::CountVoucherTypeByType - invalid input value for enum voucher_type_enum")
			return nil, errors.New(constants.ErrVoucherTypeNotFound)
		}

		log.Error().Err(err).Str("voucherType", voucherType).Msg("repository::CountVoucherTypeByType - error counting voucher type by type")
		return nil, err
	}

	return &res, nil
}
