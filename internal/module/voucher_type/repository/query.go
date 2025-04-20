package repository

const (
	queryCountVoucherTypeByType = `
		SELECT 
			id, type, COUNT(id) as count 
		FROM voucher_types 
		WHERE type = ? 
		GROUP BY id, type
	`
)
