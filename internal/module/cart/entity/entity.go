package entity

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID               int       `db:"id"`
	UserID           uuid.UUID `db:"user_id"`
	ProductID        int       `db:"product_id"`
	ProductVariantID int       `db:"product_variant_id"`
	Quantity         int       `db:"quantity"`
	CreatedAt        time.Time `db:"created_at"`
}
