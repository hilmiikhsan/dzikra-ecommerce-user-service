package service

import (
	applicationPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/list_application/ports"
	"github.com/jmoiron/sqlx"
)

var _ applicationPorts.ApplicationService = &applicationService{}

type applicationService struct {
	db                    *sqlx.DB
	applicationRepository applicationPorts.ApplicationRepository
}

func NewApplicationService(
	db *sqlx.DB,
	applicationRepository applicationPorts.ApplicationRepository,
) *applicationService {
	return &applicationService{
		db:                    db,
		applicationRepository: applicationRepository,
	}
}
