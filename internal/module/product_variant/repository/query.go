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
)
