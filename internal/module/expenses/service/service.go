package service

import (
	"context"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/expenses/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/expenses/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *expensesService) CreateExpenses(ctx context.Context, req *dto.CreateOrUpdateExpensesRequest) (*dto.CreateOrUpdateExpensesResponse, error) {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateExpenses - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::CreateExpenses - Failed to rollback transaction")
			}
		}
	}()

	date, err := utils.ParseTime(req.Date)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateExpenses - error parsing date")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrInvalidDate))
	}

	// Insert new expenses
	expensesResult, err := s.expensesRepository.InsertNewExpenses(ctx, tx, &entity.Expenses{
		Name: req.Name,
		Cost: req.Cost,
		Date: date,
	})
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateExpenses - Failed to insert new expenses")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::CreateExpenses - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrUpdateExpensesResponse{
		ID:   expensesResult.ID,
		Name: expensesResult.Name,
		Cost: expensesResult.Cost,
		Date: utils.FormatTime(expensesResult.Date),
	}, nil
}

func (s *expensesService) GetListExpenses(ctx context.Context, page, limit int, search string) (*dto.GetListExpensesResponse, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list banner
	expenses, total, err := s.expensesRepository.FindListExpenses(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListExpenses - error getting list banner")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if expenses is nil
	if expenses == nil {
		expenses = []dto.GetListExpenses{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create map response
	response := dto.GetListExpensesResponse{
		Expenses:    expenses,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    perPage,
		TotalData:   total,
	}

	// return response
	return &response, nil
}

func (s *expensesService) UpdateExpenses(ctx context.Context, req *dto.CreateOrUpdateExpensesRequest, id int) (*dto.CreateOrUpdateExpensesResponse, error) {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::UpdateExpenses - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::UpdateExpenses - Failed to rollback transaction")
			}
		}
	}()

	date, err := utils.ParseTime(req.Date)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateExpenses - error parsing date")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrInvalidDate))
	}

	// Update expenses
	expensesResult, err := s.expensesRepository.UpdateExpenses(ctx, tx, &entity.Expenses{
		ID:   id,
		Name: req.Name,
		Cost: req.Cost,
		Date: date,
	})
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrExpensesNotFound) {
			log.Error().Err(err).Msg("service::UpdateExpenses - expenses not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrExpensesNotFound))
		}

		log.Error().Err(err).Any("payload", req).Msg("service::UpdateExpenses - Failed to update expenses")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::UpdateExpenses - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrUpdateExpensesResponse{
		ID:   expensesResult.ID,
		Name: expensesResult.Name,
		Cost: expensesResult.Cost,
		Date: utils.FormatTime(expensesResult.Date),
	}, nil
}

func (s *expensesService) RemoveExpenses(ctx context.Context, id int) error {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::RemoveExpenses - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::RemoveExpenses - Failed to rollback transaction")
			}
		}
	}()

	// soft delete expenses
	if err := s.expensesRepository.SoftDeleteExpensesByID(ctx, tx, id); err != nil {
		if strings.Contains(err.Error(), constants.ErrExpensesNotFound) {
			log.Error().Err(err).Msg("service::RemoveExpenses - expenses not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrExpensesNotFound))
		}

		log.Error().Err(err).Msg("service::RemoveExpenses - Failed to soft delete expenses")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::RemoveExpenses - failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
