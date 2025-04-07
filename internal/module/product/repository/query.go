package repository

const (
	queryCountProductByName = `
		SELECT COUNT(id) FROM products WHERE name = ? AND deleted_at IS NULL
	`

	queryInsertNewProduct = `
		INSERT INTO products
		(
			name,
			real_price,
			discount_price,
			capital_price,
			description,
			specification,
			stock,
			weight,
			variant_name,
			product_category_id,
			product_sub_category_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING 
			id, 
			name, 
			real_price, 
			discount_price, 
			capital_price,
			description, 
			specification, 
			stock, 
			weight, 
			variant_name, 
			product_category_id, 
			product_sub_category_id
	`
)
