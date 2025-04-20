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
)
