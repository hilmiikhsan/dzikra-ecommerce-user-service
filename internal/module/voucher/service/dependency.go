package service

import (
	voucherPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/ports"
	voucherTypePorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher_type/ports"
	voucherUsagePorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher_usage/ports"
	"github.com/jmoiron/sqlx"
)

var _ voucherPorts.VoucherService = &voucherService{}

type voucherService struct {
	db                     *sqlx.DB
	voucherRepository      voucherPorts.VoucherRepository
	voucherTypeRepository  voucherTypePorts.VoucherTypeRepository
	voucherUsageRepository voucherUsagePorts.VoucherUsageRepository
}

func NewVoucherService(
	db *sqlx.DB,
	voucherRepository voucherPorts.VoucherRepository,
	voucherTypeRepository voucherTypePorts.VoucherTypeRepository,
	voucherUsageRepository voucherUsagePorts.VoucherUsageRepository,
) *voucherService {
	return &voucherService{
		db:                     db,
		voucherRepository:      voucherRepository,
		voucherTypeRepository:  voucherTypeRepository,
		voucherUsageRepository: voucherUsageRepository,
	}
}
