package repository

const (
	queryInsertNewProductSubCategory = `
		INSERT INTO product_sub_categories 
		(
			name,
			product_category_id
		) VALUES (?, ?) RETURNING id, name, product_category_id
	`
)
