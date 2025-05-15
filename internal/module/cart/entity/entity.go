package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID                          int             `db:"cart_id"`
	UserID                      uuid.UUID       `db:"user_id"`
	Quantity                    int             `db:"quantity"`
	ProductID                   int             `db:"product_id"`
	ProductVariantID            int             `db:"product_variant_id"`
	ProductName                 string          `db:"product_name"`
	ProductRealPrice            string          `db:"product_real_price"`
	ProductDiscountPrice        string          `db:"product_discount_price"`
	ProductStock                int             `db:"product_stock"`
	ProductWeight               float64         `db:"product_weight"`
	ProductVariantName          string          `db:"product_variant_name"`
	ProductVariantSubName       sql.NullString  `db:"product_variant_sub_name"`
	ProductVariantRealPrice     sql.NullString  `db:"product_variant_real_price"`
	ProductVariantDiscountPrice sql.NullString  `db:"product_variant_discount_price"`
	ProductVariantStock         sql.NullInt64   `db:"product_variant_stock"`
	ProductVariantWeight        sql.NullFloat64 `db:"product_variant_weight"`
	ProductGroceryID            sql.NullInt64   `db:"product_grocery_id"`
	ProductGroceryMinBuy        sql.NullInt64   `db:"product_grocery_min_buy"`
	ProductGroceryDiscount      sql.NullInt64   `db:"product_grocery_discount"`
	ProductGroceryProductID     sql.NullInt64   `db:"product_grocery_product_id"`
	ProductImageID              sql.NullInt64   `db:"product_image_id"`
	ProductImageURL             sql.NullString  `db:"product_image_url"`
	ProductImageSort            sql.NullInt64   `db:"product_image_sort"`
	ProductImageProductID       sql.NullInt64   `db:"product_image_product_id"`
	CreatedAt                   time.Time       `db:"created_at"`
}
