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

	queryFindListProductSubCategory = `
		SELECT
			psc.id,
			psc.name,
			pc.id AS category_id,
			pc.name AS category_name
		FROM product_sub_categories psc
		JOIN product_categories pc ON psc.product_category_id = pc.id
		WHERE psc.deleted_at IS NULL
		AND pc.deleted_at IS NULL
		AND pc.id = ?
		AND psc.name ILIKE '%' || ? || '%' AND psc.deleted_at IS NULL
		ORDER BY psc.name ASC
		LIMIT ? OFFSET ?
	`

	queryCountListProductSubCategory = `
		SELECT COUNT(*)
		FROM product_sub_categories psc
		JOIN product_categories pc ON psc.product_category_id = pc.id
		WHERE psc.name ILIKE '%' || ? || '%' AND psc.deleted_at IS NULL
	`

	querySoftDeleteProductSubCategory = `
		UPDATE product_sub_categories
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	queryCountProductSubCategoryByID = `
		SELECT COUNT(id)
		FROM product_sub_categories
		WHERE id = ? AND deleted_at IS NULL
	`
)
