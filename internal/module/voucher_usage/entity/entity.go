package entity

import (
	"time"

	"github.com/google/uuid"
)

type VoucherUsage struct {
	ID        int       `db:"id"`
	IsUse     bool      `db:"is_use"`
	VoucherID int       `db:"voucher_id"`
	UserID    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}
