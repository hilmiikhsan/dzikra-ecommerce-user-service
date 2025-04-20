package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/entity"
)

type VoucherRepository interface {
	InsertNewVoucher(ctx context.Context, data *entity.Voucher) (*entity.Voucher, error)
	FindListVoucher(ctx context.Context, limit, offset int, search string) ([]dto.GetListVoucher, int, error)
}

type VoucherService interface {
	CreateVoucher(ctx context.Context, req *dto.CreateVoucherRequest) (*dto.CreateVoucherResponse, error)
	GetListVoucher(ctx context.Context, page, limit int, search string) (*dto.GetListVoucherResponse, error)
}
