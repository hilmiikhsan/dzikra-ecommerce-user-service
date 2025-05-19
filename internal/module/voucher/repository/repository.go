package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/dto"
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

func (r *voucherRepository) FindListVoucher(ctx context.Context, limit, offset int, search string) ([]dto.GetListVoucher, int, error) {
	var responses []entity.Voucher

	if err := r.db.SelectContext(ctx, &responses, r.db.Rebind(queryFindListVoucher), search, limit, offset); err != nil {
		log.Error().Err(err).Msg("repository::FindListVoucher - error executing query")
		return nil, 0, err
	}

	var total int

	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountListVoucher), search); err != nil {
		log.Error().Err(err).Msg("repository::FindListVoucher - error counting voucher")
		return nil, 0, err
	}

	vouchers := make([]dto.GetListVoucher, 0, len(responses))
	for _, v := range responses {
		vouchers = append(vouchers, dto.GetListVoucher{
			ID:            v.ID,
			Name:          v.Name,
			VoucherQuota:  v.VoucherQuota,
			CreatedAt:     utils.FormatTime(v.CreatedAt),
			StartAt:       utils.FormatTime(v.StartAt),
			EndAt:         utils.FormatTime(v.EndAt),
			Code:          v.Code,
			Discount:      v.Discount,
			VoucherTypeID: v.VoucherType,
			VoucherUse:    v.VoucherUse,
		})
	}

	return vouchers, total, nil
}

func (r *voucherRepository) UpdateVoucher(ctx context.Context, data *entity.Voucher) (*entity.Voucher, error) {
	var res = new(entity.Voucher)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryUpdateVoucher),
		data.Name,
		data.VoucherQuota,
		data.Code,
		data.Discount,
		data.StartAt,
		data.EndAt,
		data.VoucherTypeID,
		data.ID,
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
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::UpdateVoucher - voucher with id %d not found", data.ID)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrVoucherNotFound)
		}

		uniqueConstraints := map[string]string{
			"vouchers_code_key": constants.ErrVoucherCodeAlreadyRegistered,
		}

		val, handleErr := utils.HandleInsertUniqueError(err, data, uniqueConstraints)
		if handleErr != nil {
			log.Error().Err(handleErr).Any("payload", data).Msg("repository::UpdateVoucher - unique violation")
			return nil, handleErr
		}

		if voucher, ok := val.(*entity.Voucher); ok {
			log.Error().Err(err).Any("payload", data).Msg("repository::UpdateVoucher - unique violation")
			return voucher, nil
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateVoucher - error updating voucher")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}

func (r *voucherRepository) SoftDeleteVoucherByID(ctx context.Context, tx *sqlx.Tx, id int) error {
	result, err := tx.ExecContext(ctx, r.db.Rebind(querySoftDeleteVoucherByID), id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("repository::SoftDeleteVoucherByID - Failed to soft delete voucher")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteVoucherByID - Failed to fetch rows affected")
		return err
	}

	if rowsAffected == 0 {
		errNotFound := errors.New(constants.ErrVoucherNotFound)
		log.Error().Err(errNotFound).Int("id", id).Msg("repository::SoftDeleteVoucherByID - Voucher not found")
		return errNotFound
	}

	return nil
}

func (r *voucherRepository) FindVoucherByCode(ctx context.Context, code string) (*entity.Voucher, error) {
	var res = new(entity.Voucher)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryFindVoucherByCode), code).Scan(
		&res.ID,
		&res.Name,
		&res.VoucherQuota,
		&res.Code,
		&res.Discount,
		&res.StartAt,
		&res.EndAt,
		&res.VoucherTypeID,
		&res.CreatedAt,
		&res.VoucherType,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Msg("repository::FindVoucherByCode - voucher not found")
			return nil, errors.New(constants.ErrVoucherNotFound)
		}

		log.Error().Err(err).Msg("repository::FindVoucherByCode - error finding voucher by code")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}

func (r *voucherRepository) FindVoucherByID(ctx context.Context, id int) (*entity.Voucher, error) {
	var res = new(entity.Voucher)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryFindVoucherByID), id).Scan(
		&res.ID,
		&res.Name,
		&res.VoucherQuota,
		&res.Code,
		&res.Discount,
		&res.StartAt,
		&res.EndAt,
		&res.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Msg("repository::FindVoucherByID - voucher not found")
			return nil, errors.New(constants.ErrVoucherNotFound)
		}

		log.Error().Err(err).Msg("repository::FindVoucherByID - error finding voucher by id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}
