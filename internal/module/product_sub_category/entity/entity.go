package entity

type ProductSubCategory struct {
	ID                int    `db:"id"`
	Name              string `db:"name"`
	ProductCategoryID int    `db:"product_category_id"`
	CategoryID        int    `db:"category_id"`
	CategoryName      string `db:"category_name"`
}
