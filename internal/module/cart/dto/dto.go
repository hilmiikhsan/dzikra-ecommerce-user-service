package dto

import (
	productGrocery "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_grocery/dto"
	productImage "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_image/dto"
)

type AddOrUpdateCartItemRequest struct {
	ProductID        int    `json:"product_id" validate:"required,numeric,non_zero_integer,gt=0"`
	ProductVariantID int    `json:"product_variant_id" validate:"required,numeric,non_zero_integer,gt=0"`
	Quantity         int    `json:"quantity" validate:"required,numeric,non_zero_integer,gt=0"`
	UserID           string `json:"user_id" validate:"required,xss_safe"`
}

type AddOrUpdateCartItemResponse struct {
	ID               int    `json:"id"`
	UserID           string `json:"user_id"`
	ProductID        int    `json:"product_id"`
	ProductVariantID int    `json:"product_variant_id"`
	Quantity         int    `json:"quantity"`
	CreatedAt        string `json:"created_at"`
}

type GetListCartResponse struct {
	ID                          int                             `json:"id"`
	Quantity                    int                             `json:"quantity"`
	ProductID                   int                             `json:"product_id"`
	ProductVariantID            int                             `json:"product_variant_id"`
	ProductName                 string                          `json:"product_name"`
	ProductRealPrice            string                          `json:"product_real_price"`
	ProductDiscountPrice        string                          `json:"product_discount_price"`
	ProductStock                int                             `json:"product_stock"`
	ProductWeight               float64                         `json:"product_weight"`
	ProductVariantWeight        float64                         `json:"product_variant_weight"`
	ProductVariantName          string                          `json:"product_variant_name"`
	ProductGrocery              []productGrocery.ProductGrocery `json:"product_grocery"`
	ProductVariantSubName       string                          `json:"product_variant_sub_name"`
	ProductVariantRealPrice     string                          `json:"product_variant_real_price"`
	ProductVariantDiscountPrice string                          `json:"product_variant_discount_price"`
	ProductVariantStock         int                             `json:"product_variant_stock"`
	ProductImage                []productImage.ProductImage     `json:"product_image"`
}
