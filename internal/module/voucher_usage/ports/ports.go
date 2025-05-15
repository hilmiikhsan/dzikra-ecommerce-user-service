package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher_usage/entity"
	"github.com/jmoiron/sqlx"
)

type VoucherUsageRepository interface {
	SoftDeleteVoucherUsageByVoucherID(ctx context.Context, tx *sqlx.Tx, voucherID int) error
	FindVoucherUsageByVoucherIdAndUserId(ctx context.Context, voucherID int, userID string) (*entity.VoucherUsage, error)
	InsertNewVoucherUsage(ctx context.Context, voucherID int, userID string) error
}
