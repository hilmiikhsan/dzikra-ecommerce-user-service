package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/entity"
	"github.com/jmoiron/sqlx"
)

type VoucherRepository interface {
	InsertNewVoucher(ctx context.Context, data *entity.Voucher) (*entity.Voucher, error)
	FindListVoucher(ctx context.Context, limit, offset int, search string) ([]dto.GetListVoucher, int, error)
	UpdateVoucher(ctx context.Context, data *entity.Voucher) (*entity.Voucher, error)
	SoftDeleteVoucherByID(ctx context.Context, tx *sqlx.Tx, id int) error
	FindVoucherByCode(ctx context.Context, code string) (*entity.Voucher, error)
	FindVoucherByID(ctx context.Context, id int) (*entity.Voucher, error)
}

type VoucherService interface {
	CreateVoucher(ctx context.Context, req *dto.CreateOrUpdateVoucherRequest) (*dto.CreateOrUpdateVoucherResponse, error)
	GetListVoucher(ctx context.Context, page, limit int, search string) (*dto.GetListVoucherResponse, error)
	UpdateVoucher(ctx context.Context, id int, req *dto.CreateOrUpdateVoucherRequest) (*dto.CreateOrUpdateVoucherResponse, error)
	RemoveVoucher(ctx context.Context, id int) error
	VoucherUse(ctx context.Context, req *dto.VoucherUseRequest, userID string) (*dto.VoucherUseResponse, error)
}
