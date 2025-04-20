package ports

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type VoucherUsageRepository interface {
	SoftDeleteVoucherUsageByVoucherID(ctx context.Context, tx *sqlx.Tx, voucherID int) error
}
