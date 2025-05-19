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

	queryFindListVoucher = `
		SELECT
			vc.id,
			vc.name,
			vc.voucher_quota,
			vc.code,
			vc.discount,
			vc.created_at,
			vc.start_at,
			vc.end_at,
			vt.type AS voucher_type,
			COALESCE((
				SELECT COUNT(*) 
				FROM voucher_usage vu
				WHERE vu.voucher_id = vc.id
				AND vu.is_use      = TRUE
				AND vu.deleted_at IS NULL
			), 0) AS voucher_use
		FROM vouchers vc
		LEFT JOIN voucher_types vt ON vc.voucher_type_id = vt.id
		WHERE 
			vc.deleted_at IS NULL AND 
			vc.name ILIKE '%' || ? || '%'
		ORDER BY vc.created_at DESC, vc.id DESC
		LIMIT ? OFFSET ?
	`

	queryCountListVoucher = `
		SELECT COUNT(*)
		FROM vouchers vc
		WHERE name ILIKE '%' || ? || '%' AND deleted_at IS NULL
	`

	queryUpdateVoucher = `
		UPDATE vouchers
		SET 
			name             = ?,
			voucher_quota    = ?,
			code             = ?,
			discount         = ?,
			start_at         = ?,
			end_at           = ?,
			voucher_type_id  = ?,
			updated_at       = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
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

	querySoftDeleteVoucherByID = `
		UPDATE vouchers 
		SET 
			deleted_at = NOW() 
		WHERE id = ? AND deleted_at IS NULL
	`

	queryFindVoucherByCode = `
		SELECT
			v.id,
			v.name,
			v.voucher_quota,
			v.code,
			v.discount,
			v.start_at,
			v.end_at,
			v.voucher_type_id,
			v.created_at,
			vt.type AS voucher_type
		FROM vouchers v
		LEFT JOIN voucher_types vt ON v.voucher_type_id = vt.id
		WHERE 
			v.deleted_at IS NULL 
			AND v.code = ?
		LIMIT 1
	`

	queryFindVoucherByID = `
		SELECT
			id,
			name,
			voucher_quota,
			code,
			discount,
			start_at,
			end_at,
			created_at
		FROM vouchers
		WHERE 
			deleted_at IS NULL 
			AND id = ?
	`
)
