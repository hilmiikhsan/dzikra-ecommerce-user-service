package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID  `db:"id"`
	Username            string     `db:"username"`
	Email               string     `db:"email"`
	Password            string     `db:"password"`
	FullName            string     `db:"full_name"`
	EmailVerifiedAt     *time.Time `db:"email_verified_at"`
	OtpNumberVerifiedAt *time.Time `db:"otp_number_verified_at"`
	LastLoginAt         *time.Time `db:"last_login_at"`
	CreatedAt           time.Time  `db:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at"`
	DeletedAt           *time.Time `db:"deleted_at"`
}

type ListUserRow struct {
	ID             string `db:"id"`
	Email          string `db:"email"`
	FullName       string `db:"full_name"`
	PhoneNumber    string `db:"phone_number"`
	UserRole       string `db:"user_role"`
	EmailConfirmed bool   `db:"email_confirmed"`
}
