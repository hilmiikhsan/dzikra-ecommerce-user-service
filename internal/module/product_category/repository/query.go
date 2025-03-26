package repostiory

const (
	queryFindListRole = `
		SELECT
			id,
			name
		FROM product_categories
		WHERE name ILIKE '%' || ? || '%' AND deleted_at IS NULL
		ORDER BY id
		LIMIT ? OFFSET ?
	`

	queryCountListRole = `
		SELECT COUNT(*)
		FROM product_categories
		WHERE name ILIKE '%' || ? || '%' AND deleted_at IS NULL
	`
)
