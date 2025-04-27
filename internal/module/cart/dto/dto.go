package dto

type AddCartItemRequest struct {
	ProductID        int    `json:"product_id" validate:"required,numeric,non_zero_integer,gt=0"`
	ProductVariantID int    `json:"product_variant_id" validate:"required,numeric,non_zero_integer,gt=0"`
	Quantity         int    `json:"quantity" validate:"required,numeric,non_zero_integer,gt=0"`
	UserID           string `json:"user_id" validate:"required,xss_safe"`
}

type AddCartItemResponse struct {
	ID               int    `json:"id"`
	UserID           string `json:"user_id"`
	ProductID        int    `json:"product_id"`
	ProductVariantID int    `json:"product_variant_id"`
	Quantity         int    `json:"quantity"`
	CreatedAt        string `json:"created_at"`
}
