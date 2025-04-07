package dto

type ProductImage struct {
	ID        int    `json:"id"`
	ImageURL  string `json:"image_url"`
	Position  int    `json:"position"`
	ProductID int    `json:"product_id"`
}
