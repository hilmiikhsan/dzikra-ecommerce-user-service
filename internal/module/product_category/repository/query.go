package repostiory

const (
	queryFindListProductCategory = `
		SELECT
			id,
			name
		FROM product_categories
		WHERE name ILIKE '%' || ? || '%' AND deleted_at IS NULL
		ORDER BY id
		LIMIT ? OFFSET ?
	`

	queryCountListProductCategory = `
		SELECT COUNT(*)
		FROM product_categories
		WHERE name ILIKE '%' || ? || '%' AND deleted_at IS NULL
	`

	queryInsertNewProductCategory = `
		INSERT INTO product_categories (name) VALUES (?) RETURNING id, name
	`

	queryUpdateProductCategory = `
		UPDATE product_categories
		SET 
			name = ?
		WHERE id = ? AND name <> ?
		RETURNING id, name
	`

	queryFindProductCategoryByID = `
		SELECT
			id,
			name
		FROM product_categories
		WHERE id = ? AND deleted_at IS NULL
	`

	queryDeleteProductCategoryByID = `
		DELETE FROM product_categories
		WHERE id = ? AND deleted_at IS NULL
	`

	queryCountProductCategoryByID = `
		SELECT COUNT(id)
		FROM product_categories
		WHERE id = ? AND deleted_at IS NULL
	`
)
