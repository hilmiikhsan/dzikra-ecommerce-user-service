package entity

import "time"

type Voucher struct {
	ID            int       `db:"id"`
	Name          string    `db:"name"`
	VoucherQuota  int       `db:"voucher_quota"`
	Code          string    `db:"code"`
	Discount      int       `db:"discount"`
	StartAt       time.Time `db:"start_at"`
	EndAt         time.Time `db:"end_at"`
	VoucherUse    int       `db:"voucher_use"`
	VoucherTypeID int       `db:"voucher_type_id"`
	VoucherType   string    `db:"voucher_type"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	DeletedAt     time.Time `db:"deleted_at"`
}
