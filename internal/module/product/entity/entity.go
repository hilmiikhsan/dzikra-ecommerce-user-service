package entity

type Product struct {
	ID                   int     `db:"id"`
	Name                 string  `db:"name"`
	RealPrice            int     `db:"real_price"`
	DiscountPrice        int     `db:"discount_price"`
	CapitalPrice         int     `db:"capital_price"`
	Description          string  `db:"description"`
	Spesification        string  `db:"specification"`
	Stock                int     `db:"stock"`
	Weight               float64 `db:"weight"`
	VariantName          string  `db:"variant_name"`
	ProductCategoryID    int     `db:"product_category_id"`
	ProductSubCategoryID int     `db:"product_sub_category_id"`
}
