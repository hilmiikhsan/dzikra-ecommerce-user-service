package dto

type CreateVoucherRequest struct {
	Name         string `json:"name" validate:"required,min=3,max=100,xss_safe"`
	Code         string `json:"code" validate:"required,min=3,max=100,xss_safe"`
	VoucherQuota int    `json:"voucher_quota" validate:"numeric,non_zero_integer,gt=0"`
	Discount     int    `json:"discount" validate:"numeric"`
	VoucherType  string `json:"voucher_type" validate:"required,min=3,max=100,xss_safe"`
	StartAt      string `json:"start_at" validate:"required,date_format"`
	EndAt        string `json:"end_at" validate:"required,date_format"`
}

type CreateVoucherResponse struct {
	Name          string `json:"name"`
	Code          string `json:"code"`
	Discount      int    `json:"discount"`
	VoucherQuota  int    `json:"voucher_quota"`
	VoucherTypeID string `json:"voucher_type_id"`
	CreatedAt     string `json:"created_at"`
	StartAt       string `json:"start_at"`
	EndAt         string `json:"end_at"`
}
