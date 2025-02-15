package dto

type RegisterRequest struct {
	FullName        string `json:"full_name" validate:"required,min=2,max=100"`
	Username        string `json:"username" validate:"required,min=2,max=50"`
	Password        string `json:"password" validate:"required,strong_password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	Email           string `json:"email" validate:"required,email,email_blacklist"`
	PhoneNumber     string `json:"phone_number" validate:"required,phone"`
}

type RegisterResponse struct {
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
}
