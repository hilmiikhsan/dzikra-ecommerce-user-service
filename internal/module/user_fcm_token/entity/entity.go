package entity

import "github.com/google/uuid"

type UserFCMToken struct {
	ID         uuid.UUID `db:"id"`
	UserID     uuid.UUID `db:"user_id"`
	DeviceID   string    `db:"device_id"`
	DeviceType string    `db:"device_type"`
	FcmToken   string    `db:"fcm_token"`
}
