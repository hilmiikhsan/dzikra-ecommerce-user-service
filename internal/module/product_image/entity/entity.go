package entity

type ProductImage struct {
	ID        int    `db:"id"`
	ImageURL  string `db:"image_url"`
	ProductID int    `db:"product_id"`
	Sort      int    `db:"sort"`
}
