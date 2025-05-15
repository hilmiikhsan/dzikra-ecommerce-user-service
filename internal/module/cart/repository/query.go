package repository

const (
	queryInsertNewCart = `
	INSERT INTO carts 
	(
		user_id,
		product_id,
		product_variant_id,
		quantity
	) VALUES (?, ?, ?, ?) 
	 RETURNING 
	 	id, 
		user_id, 
		product_id, 
		product_variant_id, 
		quantity, 
		created_at
	`

	queryFindListCartByUserID = `
		SELECT
			c.id                          AS cart_id,
			c.quantity,
			c.product_id,
			c.product_variant_id,
			p.name                        AS product_name,
			p.real_price                  AS product_real_price,
			p.discount_price              AS product_discount_price,
			p.stock                       AS product_stock,
			p.weight 				   	  AS product_weight,
			p.variant_name                AS product_variant_name,
			pv.variant_sub_name           AS product_variant_sub_name,
			pv.real_price                 AS product_variant_real_price,
			pv.discount_price             AS product_variant_discount_price,
			pv.variant_stock              AS product_variant_stock,
			pv.variant_weight			  AS product_variant_weight,
			pg.id                         AS product_grocery_id,
			pg.min_buy                    AS product_grocery_min_buy,
			pg.discount                   AS product_grocery_discount,
			pg.product_id                 AS product_grocery_product_id,
			pi.id                         AS product_image_id,
			pi.image_url                  AS product_image_url,
			pi.sort                       AS product_image_sort,
			pi.product_id                 AS product_image_product_id
		FROM carts c
		JOIN products p            ON c.product_id = p.id
		JOIN product_variants pv   ON c.product_variant_id = pv.id
		LEFT JOIN product_groceries pg ON pg.product_id = p.id AND pg.deleted_at IS NULL
		LEFT JOIN product_images pi    ON pi.product_id = p.id AND pi.deleted_at IS NULL
		WHERE c.user_id = $1 AND c.deleted_at IS NULL
		ORDER BY c.id, pi.sort
	`

	queryUpdateCart = `
		UPDATE carts
		SET
			user_id = ?,
			product_id = ?,
			product_variant_id = ?,
			quantity = ?
		WHERE id = ? AND deleted_at IS NULL
		RETURNING
			id, 
			user_id, 
			product_id, 
			product_variant_id, 
			quantity, 
			created_at
	`

	queryDeleteCartByID = `
		DELETE FROM carts WHERE id = ? AND deleted_at IS NULL
	`
)
