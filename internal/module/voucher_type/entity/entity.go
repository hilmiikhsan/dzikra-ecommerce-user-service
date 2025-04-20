package entity

type VoucherType struct {
	ID    int    `db:"id"`
	Type  string `db:"type"`
	Count int    `db:"count"`
}
