package repository

const (
	queryInsertNewProductSubCategory = `
		INSERT INTO product_sub_categories 
		(
			name,
			product_category_id
		) VALUES (?, ?) RETURNING id, name, product_category_id
	`

	queryFindProductSubCategoryByID = `
		SELECT
			id,
			name,
			product_category_id
		FROM product_sub_categories
		WHERE id = ?
		AND deleted_at IS NULL
	`

	queryUpdateProductSubCategory = `
		UPDATE product_sub_categories
		SET 
			name = ?
		WHERE id = ? AND name <> ?
		RETURNING id, name, product_category_id
	`
)
