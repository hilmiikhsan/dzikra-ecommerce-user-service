package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/dashboard/dto"
)

type DashboardService interface {
	GetDashboard(ctx context.Context, startDate, endDate string) (*dto.GetDashboardResponse, error)
}
