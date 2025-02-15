package entity

import "github.com/google/uuid"

type Role struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
