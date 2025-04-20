package repository

const (
	queryInsertNewVoucher = `
		INSERT INTO vouchers 
		(
			name, 
			voucher_quota, 
			code, 
			discount, 
			start_at, 
			end_at, 
			voucher_type_id
		) VALUES (?, ?, ?, ?, ?, ?, ?) 
		 RETURNING 
		 	id, 
			name, 
			voucher_quota, 
			code, 
			discount, 
			start_at, 
			end_at, 
			voucher_type_id, 
			created_at
	`
)
