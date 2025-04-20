package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher_type/entity"
)

type VoucherTypeRepository interface {
	CountVoucherTypeByType(ctx context.Context, voucherType string) (*entity.VoucherType, error)
}
