package service

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/proto/order"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/dashboard/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *dashboardService) GetDashboard(ctx context.Context, startDate, endDate string) (*dto.GetDashboardResponse, error) {
	start, err := utils.ParseDateToUTC(startDate)
	if err != nil {
		log.Error().Err(err).Msg("service::GetDashboard - error parsing start date")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrInvalidStartDate))
	}

	end, err := utils.ParseEndDateToUTC(endDate)
	if err != nil {
		log.Error().Err(err).Msg("service::GetDashboard - error parsing end date")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrInvalidEndDate))
	}

	totalSumExpenses, err := s.expensesRepository.FindTotalSumExpenses(ctx, start, end)
	if err != nil {
		log.Error().Err(err).Msg("service::GetDashboard - error getting total sum expenses")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	calculateSumOrderResult, err := s.externalOrder.CalculateTotalSummary(ctx, &order.CalculateTotalSummaryRequest{
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		log.Error().Err(err).Msg("service::GetDashboard - error calculating total summary")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	profitLoss := calculateSumOrderResult.NetSales - totalSumExpenses

	return &dto.GetDashboardResponse{
		TotalAmount:         int(calculateSumOrderResult.TotalAmount),
		TotalExpenses:       int(totalSumExpenses),
		TotalTransaction:    int(calculateSumOrderResult.TotalTransaction),
		TotalSellingProduct: int(calculateSumOrderResult.TotalSellingProduct),
		TotalCapital:        int(calculateSumOrderResult.TotalCapital),
		Netsales:            int(calculateSumOrderResult.NetSales),
		ProfitLoss:          int(profitLoss),
	}, nil
}
