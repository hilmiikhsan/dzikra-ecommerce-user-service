package repository

const (
	queryInsertNewBanner = `
		INSERT INTO banners 
		(
			image_url, 
			description
		) VALUES (?, ?) RETURNING id, image_url, description
	`
)
