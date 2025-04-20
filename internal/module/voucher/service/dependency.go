package service

import (
	voucherPorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/ports"
	voucherTypePorts "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher_type/ports"
)

var _ voucherPorts.VoucherService = &voucherService{}

type voucherService struct {
	voucherRepository     voucherPorts.VoucherRepository
	voucherTypeRepository voucherTypePorts.VoucherTypeRepository
}

func NewVoucherService(
	voucherRepository voucherPorts.VoucherRepository,
	voucherTypeRepository voucherTypePorts.VoucherTypeRepository,
) *voucherService {
	return &voucherService{
		voucherRepository:     voucherRepository,
		voucherTypeRepository: voucherTypeRepository,
	}
}
