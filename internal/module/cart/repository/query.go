package repository

const (
	queryInsertNewCart = `
	INSERT INTO carts 
	(
		user_id,
		product_id,
		product_variant_id,
		quantity
	) VALUES (?, ?, ?, ?) 
	 RETURNING 
	 	id, 
		user_id, 
		product_id, 
		product_variant_id, 
		quantity, 
		created_at
	`
)
