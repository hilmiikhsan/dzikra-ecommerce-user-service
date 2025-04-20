package repository

const (
	queryInsertNewBanner = `
		INSERT INTO banners 
		(
			image_url, 
			description
		) VALUES (?, ?) RETURNING id, image_url, description
	`

	queryFindListBanner = `
		SELECT
			id,
			image_url,
			description
		FROM banners
		WHERE
			deleted_at IS NULL AND
			description ILIKE '%' || ? || '%'
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?
	`

	queryCountListBanner = `
		SELECT COUNT(*)
		FROM banners
		WHERE description ILIKE '%' || ? || '%' AND deleted_at IS NULL
	`

	queryUpdateBanner = `
		UPDATE banners
		SET
			image_url = ?,
			description = ?
		WHERE id = ? AND deleted_at IS NULL
		RETURNING id, image_url, description
	`

	queryFindBannerByID = `
		SELECT
			id,
			image_url,
			description
		FROM banners
		WHERE id = ? AND deleted_at IS NULL
	`

	querySoftDeleteBannerByID = `
		UPDATE banners
		SET 
			deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
)
