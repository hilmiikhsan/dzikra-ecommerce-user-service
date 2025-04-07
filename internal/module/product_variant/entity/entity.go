package entity

type ProductVariant struct {
	ID             int     `db:"id"`
	VariantSubName string  `db:"variant_sub_name"`
	VariantStock   int     `db:"variant_stock"`
	VariantWeight  float64 `db:"variant_weight"`
	CapitalPrice   int     `db:"capital_price"`
	RealPrice      int     `db:"real_price"`
	DiscountPrice  int     `db:"discount_price"`
	ProductID      int     `db:"product_id"`
}
