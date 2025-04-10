package entity

import (
	"database/sql"

	productGrocery "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/entity"
	productImage "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/entity"
	productVariant "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_variant/entity"
)

type Product struct {
	ID                          int                             `db:"id"`
	Name                        string                          `db:"name"`
	RealPrice                   int                             `db:"real_price"`
	DiscountPrice               int                             `db:"discount_price"`
	CapitalPrice                int                             `db:"capital_price"`
	Description                 string                          `db:"description"`
	Spesification               string                          `db:"specification"`
	Stock                       int                             `db:"stock"`
	Weight                      float64                         `db:"weight"`
	VariantName                 string                          `db:"variant_name"`
	ProductCategoryID           int                             `db:"product_category_id"`
	ProductCategoryName         string                          `db:"product_category_name"`
	ProductSubID                int                             `db:"product_sub_id"`
	ProductSubCategoryName      string                          `db:"product_sub_category_name"`
	ProductSubCategoryID        int                             `db:"product_sub_category_id"`
	ProductVariantID            sql.NullInt64                   `db:"product_variant_id"`
	ProductVariantSubName       sql.NullString                  `db:"product_variant_sub_name"`
	ProductVariantStock         sql.NullInt64                   `db:"product_variant_stock"`
	ProductVariantWeight        sql.NullFloat64                 `db:"product_variant_weight"`
	ProductVariantCapitalPrice  sql.NullInt64                   `db:"product_variant_capital_price"`
	ProductVariantRealPrice     sql.NullInt64                   `db:"product_variant_real_price"`
	ProductVariantDiscountPrice sql.NullInt64                   `db:"product_variant_discount_price"`
	ProductVariantProductID     sql.NullInt64                   `db:"product_variant_product_id"`
	ProductGroceryID            sql.NullInt64                   `db:"product_grocery_id"`
	ProductGroceryMinBuy        sql.NullInt64                   `db:"product_grocery_min_buy"`
	ProductGroceryDiscount      sql.NullInt64                   `db:"product_grocery_discount"`
	ProductGroceryProductID     sql.NullInt64                   `db:"product_grocery_product_id"`
	ProductImageID              sql.NullInt64                   `db:"product_image_id"`
	ProductImageURL             sql.NullString                  `db:"product_image_url"`
	ProductImageSort            sql.NullInt64                   `db:"product_image_sort"`
	ProductImageProductID       sql.NullInt64                   `db:"product_image_product_id"`
	ProductVariant              []productVariant.ProductVariant `db:"-"`
	ProductGrocery              []productGrocery.ProductGrocery `db:"-"`
	ProductImage                []productImage.ProductImage     `db:"-"`
}
