package repository

const (
	querySoftDeleteVoucherUsageByVoucherID = `
		UPDATE voucher_usage 
		SET 
			deleted_at = NOW() 
		WHERE voucher_id = ? AND deleted_at IS NULL	
	`

	queryFindVoucherUsageByVoucherIdAndUserId = `
		SELECT
			id,
			is_use,
			voucher_id,
			user_id,
			created_at
		FROM voucher_usage
		WHERE 
			voucher_id = ? 
			AND user_id = ? 
			AND deleted_at IS NULL
		LIMIT 1
	`

	queryLockVoucherRow = `
		SELECT 
			voucher_quota
        FROM vouchers
        WHERE id = $1
        FOR UPDATE
	`

	queryDecrementVoucher = `
		UPDATE vouchers
        SET 
			voucher_quota = voucher_quota - 1
        WHERE id = $1
	`

	queryInsertNewVoucherUsage = `
		INSERT INTO voucher_usage 
		(
			is_use, 
			voucher_id, 
			user_id
		) VALUES (TRUE, $1, $2)
	`
)
