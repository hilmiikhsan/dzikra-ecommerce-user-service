package entity

import "github.com/google/uuid"

type AppPermission struct {
	ID uuid.UUID `db:"id"`
}
