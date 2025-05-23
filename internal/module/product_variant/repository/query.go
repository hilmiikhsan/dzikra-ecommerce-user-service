package repository

const (
	queryInsertNewProductVariant = `
		INSERT INTO product_variants
		(
			variant_sub_name,
			variant_stock,
			variant_weight,
			capital_price,
			real_price,
			discount_price,
			product_id
		) VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING 
			id, 
			variant_sub_name,
			variant_stock,
			variant_weight,
			capital_price,
			real_price,
			discount_price,
			product_id
	`

	queryUpdateProductVariant = `
		UPDATE product_variants
		SET 
			variant_sub_name = ?,
			variant_stock = ?,
			variant_weight = ?,
			capital_price = ?,
			real_price = ?,
			discount_price = ?
		WHERE id = ? AND product_id = ?
		RETURNING 
			id,
			variant_sub_name,
			variant_stock,
			variant_weight,
			capital_price,
			real_price,
			discount_price,
			product_id
	`

	queryDeleteProductVariant = `
		DELETE FROM product_variants
		WHERE id = ? AND product_id = ?
	`

	querySoftDeleteProductVariantsByProductID = `
		UPDATE product_variants 
		SET 
			deleted_at = NOW() 
		WHERE product_id = ? AND deleted_at IS NULL
	`

	queryCountProductVariantByIDAndProductID = `
		SELECT COUNT(id) 		
		FROM product_variants
		WHERE id = ? 
		AND product_id = ?
		AND deleted_at IS NULL
	`

	queryFindProductVariantStockByID = `
		SELECT variant_stock
		FROM product_variants
		WHERE id = ?
		AND deleted_at IS NULL
	`
)
