package dto

type ProductVariant struct {
	ID             int     `json:"id"`
	VariantSubName string  `json:"variant_sub_name"`
	VariantStock   int     `json:"variant_stock"`
	VariantWeight  float64 `json:"variant_weight"`
	CapitalPrice   int     `json:"capital_price"`
	RealPrice      int     `json:"real_price"`
	DiscountPrice  int     `json:"discount_price"`
	ProductID      int     `json:"product_id"`
}

type Variant struct {
	VariantSubName string  `json:"variant_sub_name" validate:"required,max=100,xss_safe"`
	VariantStock   int     `json:"variant_stock" validate:"numeric"`
	RealPrice      int     `json:"real_price" validate:"numeric"`
	DiscountPrice  int     `json:"discount_price" validate:"numeric"`
	CapitalPrice   int     `json:"capital_price" validate:"numeric"`
	VariantWeight  float64 `json:"variant_weight" validate:"numeric"`
}
