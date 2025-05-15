package repository

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher_usage/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher_usage/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
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

func (r *voucherUsageRepository) FindVoucherUsageByVoucherIdAndUserId(ctx context.Context, voucherID int, userID string) (*entity.VoucherUsage, error) {
	var res = new(entity.VoucherUsage)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryFindVoucherUsageByVoucherIdAndUserId), voucherID, userID).Scan(
		&res.ID,
		&res.IsUse,
		&res.VoucherID,
		&res.UserID,
		&res.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Msg("repository::FindVoucherUsageByVoucherIdAndUserId - no rows found")
			return nil, nil
		}

		log.Error().Err(err).Msg("repository::FindVoucherUsageByVoucherIdAndUserId - error finding voucher usage")
		return nil, err
	}

	return res, nil
}

func (r *voucherUsageRepository) InsertNewVoucherUsage(ctx context.Context, voucherID int, userID string) error {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Error().Err(err).Msg("repository::InsertNewVoucherUsage - error starting transaction")
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("voucher_id", voucherID).Msg("repository::InsertNewVoucherUsage - Failed to rollback transaction")
			}
		}
	}()

	var quota int
	if err := tx.GetContext(ctx,
		&quota,
		queryLockVoucherRow,
		voucherID,
	); err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("repository::InsertNewVoucherUsage - lock failed")
		return err
	}

	if quota <= 0 {
		log.Error().Msg("repository::InsertNewVoucherUsage - voucher has been run out")
		return err_msg.NewCustomErrors(http.StatusBadRequest, err_msg.WithMessage(constants.ErrVoucherHasBeenRunOut))
	}

	if _, err := tx.ExecContext(ctx,
		queryDecrementVoucher,
		voucherID,
	); err != nil {
		log.Error().Err(err).Msg("repository::InsertNewVoucherUsage - decrement failed")
		return err
	}

	if _, err := tx.ExecContext(ctx,
		queryInsertNewVoucherUsage,
		voucherID, userID,
	); err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("repository::InsertNewVoucherUsage - insert usage failed")
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("repository::InsertNewVoucherUsage - commit failed")
		return err
	}

	log.Info().Msg("repository::InsertNewVoucherUsage - success insert voucher usage")

	return nil
}
