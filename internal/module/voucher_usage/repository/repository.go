package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher_usage/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.VoucherUsageRepository = &voucherUsageRepository{}

type voucherUsageRepository struct {
	db *sqlx.DB
}

func NewVoucherUsageRepository(db *sqlx.DB) *voucherUsageRepository {
	return &voucherUsageRepository{
		db: db,
	}
}

func (r *voucherUsageRepository) SoftDeleteVoucherUsageByVoucherID(ctx context.Context, tx *sqlx.Tx, voucherID int) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(querySoftDeleteVoucherUsageByVoucherID), voucherID)
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteVoucherUsageByVoucherID - error soft deleting voucher")
		return err
	}

	return nil
}
