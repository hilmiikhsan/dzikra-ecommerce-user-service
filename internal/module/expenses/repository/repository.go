package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/expenses/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/expenses/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/expenses/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ExpensesRepository = &expensesRepository{}

type expensesRepository struct {
	db *sqlx.DB
}

func NewExpensesRepository(db *sqlx.DB) *expensesRepository {
	return &expensesRepository{
		db: db,
	}
}

func (r *expensesRepository) InsertNewExpenses(ctx context.Context, tx *sqlx.Tx, data *entity.Expenses) (*entity.Expenses, error) {
	var res = new(entity.Expenses)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewExpenses),
		data.Name,
		data.Cost,
		data.Date,
	).Scan(
		&res.ID,
		&res.Name,
		&res.Cost,
		&res.Date,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewExpenses - Failed to insert new expenses")
		return nil, err
	}

	return res, nil
}

func (r *expensesRepository) FindListExpenses(ctx context.Context, limit, offset int, search string) ([]dto.GetListExpenses, int, error) {
	var responses []entity.Expenses

	if err := r.db.SelectContext(ctx, &responses, r.db.Rebind(queryFindListExpenses), search, limit, offset); err != nil {
		log.Error().Err(err).Msg("repository::FindListExpenses - error executing query")
		return nil, 0, err
	}

	var total int

	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountFindListExpenses), search); err != nil {
		log.Error().Err(err).Msg("repository::FindListExpenses - error counting expenses")
		return nil, 0, err
	}

	expenses := make([]dto.GetListExpenses, 0, len(responses))
	for _, v := range responses {
		expenses = append(expenses, dto.GetListExpenses{
			ID:   v.ID,
			Name: v.Name,
			Cost: v.Cost,
			Date: utils.FormatTime(v.Date),
		})
	}

	return expenses, total, nil
}

func (r *expensesRepository) UpdateExpenses(ctx context.Context, tx *sqlx.Tx, data *entity.Expenses) (*entity.Expenses, error) {
	var res = new(entity.Expenses)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryUpdateExpenses),
		data.Name,
		data.Cost,
		data.Date,
		data.ID,
	).Scan(
		&res.ID,
		&res.Name,
		&res.Cost,
		&res.Date,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::UpdateExpenses - address with id %d is not found", data.ID)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrExpensesNotFound)
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateExpenses - Failed to update expenses")
		return nil, err
	}

	return res, nil
}

func (r *expensesRepository) SoftDeleteExpensesByID(ctx context.Context, tx *sqlx.Tx, id int) error {
	result, err := tx.ExecContext(ctx, r.db.Rebind(querySoftDeleteExpensesByID), id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("repository::SoftDeleteExpensesByID - Failed to soft delete expenses")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteExpensesByID - Failed to fetch rows affected")
		return err
	}

	if rowsAffected == 0 {
		errNotFound := errors.New(constants.ErrExpensesNotFound)
		log.Error().Err(errNotFound).Int("id", id).Msg("repository::SoftDeleteExpensesByID - Expenses not found")
		return errNotFound
	}

	return nil
}

func (r *expensesRepository) FindTotalSumExpenses(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	var total float64

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryFindTotalSumExpenses), startDate, endDate).Scan(&total)
	if err != nil {
		log.Error().Err(err).Msg("repository::FindTotalSumExpenses - error executing query")
		return 0, err
	}

	return total, nil
}
