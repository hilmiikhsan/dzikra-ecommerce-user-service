package repository

const (
	queryInsertNewProductImage = `
		INSERT INTO product_images
		(
			image_url,
			product_id,
			sort
		) VALUES (?, ?, ?) RETURNING id, product_id, sort
	`

	queryGetNextSort = `
		SELECT COALESCE(MAX(sort), 0) + 1 
		FROM product_images 
		WHERE product_id = ? AND deleted_at IS NULL
	`

	queryFindOrderProductImage = `
		SELECT id 
		FROM product_images 
		WHERE product_id = ? AND deleted_at IS NULL
		ORDER BY sort ASC
	`

	queryUpdateProductImageSorting = `
		UPDATE product_images
		SET 
			sort = ?
		WHERE id = ?
	`

	queryCountProductImages = `
		SELECT COUNT(id) FROM product_images
		WHERE product_id = ? AND deleted_at IS NULL
	`

	queryUpdateProductImageURL = `
		UPDATE product_images
		SET 
			image_url = ?
		WHERE id = ? AND deleted_at IS NULL
		RETURNING id, image_url, product_id, sort
	`

	queryFindProductImagesByProductID = `
		SELECT 
			id, 
			image_url, 
			product_id, 
			sort
		FROM product_images
		WHERE product_id = ? AND deleted_at IS NULL
		ORDER BY sort ASC
	`

	queryDeleteProductImage = `
		UPDATE product_images
		SET 
			deleted_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`
)
