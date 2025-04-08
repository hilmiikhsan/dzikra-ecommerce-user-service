package dto

type ProductGrocery struct {
	ID        int `json:"id"`
	MinBuy    int `json:"min_buy"`
	Discount  int `json:"discount"`
	ProductID int `json:"product_id"`
}

type GroceryPrice struct {
	ID       int `json:"id,omitempty"`
	MinBuy   int `json:"min_buy" validate:"numeric"`
	Discount int `json:"discount" validate:"numeric"`
}
