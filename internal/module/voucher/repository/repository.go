package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.VoucherRepository = &voucherRepository{}

type voucherRepository struct {
	db *sqlx.DB
}

func NewVoucherRepository(db *sqlx.DB) *voucherRepository {
	return &voucherRepository{
		db: db,
	}
}

func (r *voucherRepository) InsertNewVoucher(ctx context.Context, data *entity.Voucher) (*entity.Voucher, error) {
	var res = new(entity.Voucher)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryInsertNewVoucher),
		data.Name,
		data.VoucherQuota,
		data.Code,
		data.Discount,
		data.StartAt,
		data.EndAt,
		data.VoucherTypeID,
	).Scan(
		&res.ID,
		&res.Name,
		&res.VoucherQuota,
		&res.Code,
		&res.Discount,
		&res.StartAt,
		&res.EndAt,
		&res.VoucherTypeID,
		&res.CreatedAt,
	)
	if err != nil {
		uniqueConstraints := map[string]string{
			"vouchers_code_key": constants.ErrVoucherCodeAlreadyRegistered,
		}

		val, handleErr := utils.HandleInsertUniqueError(err, data, uniqueConstraints)
		if handleErr != nil {
			log.Error().Err(handleErr).Any("payload", data).Msg("repository::InsertNewVoucher - Failed to insert new voucher")
			return nil, handleErr
		}

		if voucher, ok := val.(*entity.Voucher); ok {
			log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewVoucher - Failed to insert new voucher")
			return voucher, nil
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewVoucher - Failed to insert new voucher")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}
