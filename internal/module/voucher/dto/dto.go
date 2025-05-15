package dto

type CreateOrUpdateVoucherRequest struct {
	Name         string `json:"name" validate:"required,min=3,max=100,xss_safe"`
	Code         string `json:"code" validate:"required,min=3,max=100,xss_safe"`
	VoucherQuota int    `json:"voucher_quota" validate:"numeric,non_zero_integer,gt=0"`
	Discount     int    `json:"discount" validate:"numeric"`
	VoucherType  string `json:"voucher_type" validate:"required,min=3,max=100,xss_safe"`
	StartAt      string `json:"start_at" validate:"required,date_format"`
	EndAt        string `json:"end_at" validate:"required,date_format"`
}

type CreateOrUpdateVoucherResponse struct {
	Name          string `json:"name"`
	Code          string `json:"code"`
	Discount      int    `json:"discount"`
	VoucherQuota  int    `json:"voucher_quota"`
	VoucherTypeID string `json:"voucher_type_id"`
	CreatedAt     string `json:"created_at"`
	StartAt       string `json:"start_at"`
	EndAt         string `json:"end_at"`
}

type GetListVoucherResponse struct {
	Voucher     []GetListVoucher `json:"voucher"`
	TotalPages  int              `json:"total_page"`
	CurrentPage int              `json:"current_page"`
	PageSize    int              `json:"page_size"`
	TotalData   int              `json:"total_data"`
}

type GetListVoucher struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	VoucherQuota  int    `json:"voucher_quota"`
	Code          string `json:"code"`
	Discount      int    `json:"discount"`
	VoucherUse    int    `json:"voucher_use"`
	VoucherTypeID string `json:"voucher_type_id"`
	CreatedAt     string `json:"created_at"`
	StartAt       string `json:"start_at"`
	EndAt         string `json:"end_at"`
}

type VoucherUseRequest struct {
	Code string `json:"code" validate:"required,min=3,max=15,xss_safe"`
}

type VoucherUseResponse struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	VoucherQuota  int    `json:"voucher_quota"`
	Code          string `json:"code"`
	Discount      int    `json:"discount"`
	VoucherTypeID string `json:"voucher_type_id"`
	CreatedAt     string `json:"created_at"`
	StartAt       string `json:"start_at"`
	EndAt         string `json:"end_at"`
}
