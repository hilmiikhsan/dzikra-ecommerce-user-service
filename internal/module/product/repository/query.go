package repository

const (
	queryInsertNewProduct = `
		INSERT INTO products
		(
			name,
			real_price,
			discount_price,
			capital_price,
			description,
			specification,
			stock,
			weight,
			variant_name,
			product_category_id,
			product_sub_category_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING 
			id, 
			name, 
			real_price, 
			discount_price, 
			capital_price,
			description, 
			specification, 
			stock, 
			weight, 
			variant_name, 
			product_category_id, 
			product_sub_category_id
	`

	queryUpdateProduct = `
		UPDATE products
		SET 
			name = ?,
			real_price = ?,
			discount_price = ?,
			capital_price = ?,
			description = ?,
			specification = ?,
			stock = ?,
			weight = ?,
			variant_name = ?,
			product_category_id = ?,
			product_sub_category_id = ?
		WHERE id = ?
		RETURNING 
			id, 
			name, 
			real_price, 
			discount_price, 
			capital_price, 
			description, 
			specification, 
			stock, 
			weight, 
			variant_name, 
			product_category_id, 
			product_sub_category_id
	`

	queryCountProductByName = `
		SELECT COUNT(id) FROM products WHERE name = ? AND deleted_at IS NULL
	`

	queryCountListProduct = `
		SELECT COUNT(*)
		FROM products p
		WHERE 
			p.deleted_at IS NULL 
			AND ($1 = '' OR p.name ILIKE '%' || $1 || '%')
			AND ($2 = 0 OR p.product_category_id = $2)
			AND ($3 = 0 OR p.product_sub_category_id = $3)

	`

	queryFindListProduct = `
		SELECT 
			p.id AS id,
			p.name AS name,
			p.description AS description,
			p.specification AS specification,
			p.real_price AS real_price,
			p.capital_price AS capital_price,
			p.discount_price AS discount_price,
			p.stock AS stock,
			p.weight AS weight,
			p.variant_name AS variant_name,
			pc.id AS product_category_id,
			pc.name AS product_category_name,
			psc.id AS product_sub_id,
			psc.name AS product_sub_category_name,
			psc.product_category_id AS product_sub_category_id,
			pv.id AS product_variant_id,
			pv.variant_sub_name AS product_variant_sub_name,
			pv.variant_stock AS product_variant_stock,
			pv.variant_weight AS product_variant_weight,
			pv.capital_price AS product_variant_capital_price,
			pv.real_price AS product_variant_real_price,
			pv.discount_price AS product_variant_discount_price,
			pv.product_id AS product_variant_product_id,
			pg.id AS product_grocery_id,
			pg.min_buy AS product_grocery_min_buy,
			pg.discount AS product_grocery_discount,
			pg.product_id AS product_grocery_product_id,
			pi.id AS product_image_id,
			pi.image_url AS product_image_url,
			pi.sort AS product_image_sort,
			pi.product_id AS product_image_product_id
		FROM products p
		JOIN product_categories pc ON p.product_category_id = pc.id
		JOIN product_sub_categories psc ON p.product_sub_category_id = psc.id
		LEFT JOIN product_variants pv ON pv.product_id = p.id
		LEFT JOIN product_groceries pg ON pg.product_id = p.id
		LEFT JOIN product_images pi ON pi.product_id = p.id AND pi.deleted_at IS NULL
		WHERE 
			p.deleted_at IS NULL  
			AND ($1 = '' OR p.name ILIKE '%' || $1 || '%')
			AND ($2 = 0 OR p.product_category_id = $2)
			AND ($3 = 0 OR p.product_sub_category_id = $3)
		ORDER BY p.created_at DESC, p.id DESC
		LIMIT $4 OFFSET $5
	`

	queryFindProductByID = `
		SELECT 
			id,
			name,
			description,
			specification,
			real_price,
			capital_price,
			discount_price,
			stock,
			weight,
			variant_name,
			product_category_id,
			product_sub_category_id
		FROM products
		WHERE id = ? AND deleted_at IS NULL
	`

	queryFindVariants = `
		SELECT 
			id,
			variant_sub_name,
			variant_stock,
			variant_weight,
			capital_price,
			real_price,
			discount_price,
			product_id
		FROM product_variants
		WHERE product_id = ?
	`

	queryFindGroceries = `
		SELECT 
			id,
			min_buy,
			discount,
			product_id
		FROM product_groceries
		WHERE product_id = ?
	`

	queryFindImages = `
		SELECT 
			id,
			image_url,
			sort,
			product_id
		FROM product_images
		WHERE product_id = ? AND deleted_at IS NULL
		ORDER BY sort
	`

	querySoftDeleteProductByID = `
		UPDATE products 
		SET 
			deleted_at = NOW() 
		WHERE id = ? AND deleted_at IS NULL
	`

	queryCountProductByID = `
		SELECT COUNT(id) 
		FROM products 
		WHERE id = ? AND deleted_at IS NULL
	`

	queryFindProductStockByID = `
		SELECT stock
		FROM products
		WHERE 
			id = ? 
			AND deleted_at IS NULL
	`
)
