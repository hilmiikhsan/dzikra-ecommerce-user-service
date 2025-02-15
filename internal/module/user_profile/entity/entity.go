package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	ID          uuid.UUID  `db:"id"`
	UserID      uuid.UUID  `db:"user_id"`
	BirthDate   *time.Time `db:"date_of_birth"`
	PhoneNumber *string    `db:"phone_number"`
	AvatarUrl   *string    `db:"avatar_url"`
}
