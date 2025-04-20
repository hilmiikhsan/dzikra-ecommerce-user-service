package entity

type Banner struct {
	ID          int    `db:"id"`
	ImageURL    string `db:"image_url"`
	Description string `db:"description"`
}
