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
)
