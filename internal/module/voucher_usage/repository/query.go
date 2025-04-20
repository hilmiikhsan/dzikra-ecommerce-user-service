package repository

const (
	querySoftDeleteVoucherUsageByVoucherID = `
		UPDATE voucher_usage 
		SET 
			deleted_at = NOW() 
		WHERE voucher_id = ? AND deleted_at IS NULL	
	`
)
