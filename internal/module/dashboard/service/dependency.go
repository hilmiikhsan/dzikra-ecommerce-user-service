package service

import (
	externalOrder "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/external/order"
	dashboardPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/dashboard/ports"
	expensesPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/expenses/ports"
	"github.com/jmoiron/sqlx"
)

var _ dashboardPorts.DashboardService = &dashboardService{}

type dashboardService struct {
	db                 *sqlx.DB
	expensesRepository expensesPorts.ExpensesRepository
	externalOrder      externalOrder.ExternalOrder
}

func NewDashboardService(
	db *sqlx.DB,
	expensesRepository expensesPorts.ExpensesRepository,
	externalOrder externalOrder.ExternalOrder,
) *dashboardService {
	return &dashboardService{
		db:                 db,
		expensesRepository: expensesRepository,
		externalOrder:      externalOrder,
	}
}
