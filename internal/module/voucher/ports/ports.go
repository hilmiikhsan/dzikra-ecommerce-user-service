package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/entity"
)

type VoucherRepository interface {
	InsertNewVoucher(ctx context.Context, data *entity.Voucher) (*entity.Voucher, error)
}

type VoucherService interface {
	CreateVoucher(ctx context.Context, req *dto.CreateVoucherRequest) (*dto.CreateVoucherResponse, error)
}
