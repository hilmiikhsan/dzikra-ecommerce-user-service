package repository

const (
	queryInsertNewProductGrocery = `
		INSERT INTO product_groceries
		(
			min_buy, 
			discount, 
			product_id
		) VALUES (?, ?, ?) RETURNING id, min_buy, discount, product_id
	`

	queryUpdateProductGrocery = `
		UPDATE product_groceries
		SET 
			min_buy = ?,
			discount = ?
		WHERE id = ? AND product_id = ?
		RETURNING 
			id,
			min_buy,
			discount,
			product_id
	`

	querySoftDeleteProductGroceriesByProductID = `
		UPDATE product_groceries 
		SET 
			deleted_at = NOW() 
		WHERE product_id = ? AND deleted_at IS NULL
	`

	queryFindProductGroceryByProductID = `
		SELECT
			id,
			min_buy,
			discount,
			product_id
		FROM product_groceries
		WHERE 
			product_id = ? 
			AND deleted_at IS NULL
	`
)
