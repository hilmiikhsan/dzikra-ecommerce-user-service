package repository

const (
	queryInsertNewAddress = `
	INSERT INTO addresses
	(
		province,
		city,
		district,
		subdistrict,
		postal_code,
		address,
		received_name,
		user_id,
		city_vendor_id,
		province_vendor_id,
		subdistrict_vendor_id
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING 
		id, 
		province, 
		city, 
		district, 
		subdistrict, 
		postal_code, 
		address, 
		received_name, 
		user_id, 
		city_vendor_id, 
		province_vendor_id, 
		subdistrict_vendor_id
	`

	queryUpdateAddress = `
		UPDATE addresses
		SET
			province = ?,
			city = ?,
			district = ?,
			subdistrict = ?,
			postal_code = ?,
			address = ?,
			received_name = ?,
			user_id = ?,
			city_vendor_id = ?,
			province_vendor_id = ?,
			subdistrict_vendor_id = ?
		WHERE id = ? AND deleted_at IS NULL
		RETURNING
			id,
			province,
			city,
			district,
			subdistrict,
			postal_code,
			address,
			received_name,
			user_id,
			city_vendor_id,
			province_vendor_id,
			subdistrict_vendor_id
	`
)
