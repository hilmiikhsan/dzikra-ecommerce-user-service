package entity

type ProductGrocery struct {
	ID        int `db:"id"`
	MinBuy    int `db:"min_buy"`
	Discount  int `db:"discount"`
	ProductID int `db:"product_id"`
}
