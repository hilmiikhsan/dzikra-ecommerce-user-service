package dto

type CreateOrderRequest struct {
	CostName       string  `json:"cost_name" validate:"required,min=3,max=30,xss_safe"`
	CostService    string  `json:"cost_service" validate:"required,min=3,max=30,xss_safe"`
	AddressID      string  `json:"address_id" validate:"required,xss_safe"`
	CallbackFinish string  `json:"callback_finish" validate:"required,url,xss_safe"`
	VoucherID      *string `json:"voucher_id" validate:"omitempty,xss_safe"`
	Notes          string  `json:"notes" validate:"required,min=2,max=100,xss_safe"`
}

type CreateOrderResponse struct {
	Order               OrderDetail `json:"order"`
	MidtransRedirectUrl string      `json:"midtrans_redirect_url"`
}

type OrderDetail struct {
	ID                  string `json:"id"`
	OrderDate           string `json:"order_date"`
	Status              string `json:"status"`
	ShippingName        string `json:"shipping_name"`
	ShippingAddress     string `json:"shipping_address"`
	ShippingPhone       string `json:"shipping_phone"`
	ShippingNumber      string `json:"shipping_number"`
	ShippingType        string `json:"shipping_type"`
	TotalWeight         int    `json:"total_weight"`
	TotalQuantity       int    `json:"total_quantity"`
	TotalShippingCost   string `json:"total_shipping_cost"`
	TotalProductAmount  string `json:"total_product_amount"`
	TotalShippingAmount string `json:"total_shipping_amount"`
	TotalAmount         string `json:"total_amount"`
	VoucherDiscount     int    `json:"voucher_discount"`
	VoucherID           string `json:"voucher_id"`
	CostName            string `json:"cost_name"`
	CostService         string `json:"cost_service"`
	AddressID           int    `json:"address_id"`
	UserID              string `json:"user_id"`
	Notes               string `json:"notes"`
}

type GetListOrderResponse struct {
	Orders      []GetListOrder `json:"orders"`
	TotalPages  int            `json:"total_page"`
	CurrentPage int            `json:"current_page"`
	PageSize    int            `json:"page_size"`
	TotalData   int            `json:"total_data"`
}

type GetListOrder struct {
	ID                  string      `json:"id"`
	OrderDate           string      `json:"order_date"`
	Status              string      `json:"status"`
	TotalQuantity       int         `json:"total_quantity"`
	TotalAmount         string      `json:"total_amount"`
	ShippingNumber      string      `json:"shipping_number"`
	TotalShippingAmount string      `json:"total_shipping_amount"`
	CostName            string      `json:"cost_name"`
	CostService         string      `json:"cost_service"`
	VoucherID           *string     `json:"voucher_id"`
	VoucherDiscount     int         `json:"voucher_disc"`
	UserID              string      `json:"user_id"`
	Notes               string      `json:"notes"`
	SubTotal            string      `json:"sub_total"`
	Address             Address     `json:"address"`
	OrderItems          []OrderItem `json:"order_item"`
	Payment             Payment     `json:"payment"`
}

type OrderItem struct {
	ProductID             int            `json:"product_id"`
	ProductName           string         `json:"product_name"`
	ProductVariantSubName string         `json:"product_variant_sub_name"`
	ProductVariant        string         `json:"product_variant"`
	TotalAmount           string         `json:"total_amount"`
	ProductDisc           *string        `json:"product_disc"`
	Quantity              int            `json:"quantity"`
	FixPricePerItem       string         `json:"fix_price_per_item"`
	ProductImages         []ProductImage `json:"product_image"`
}

type Payment struct {
	RedirectURL string `json:"redirect_url"`
}

type ProductImage struct {
	ID        int    `json:"id"`
	ImageURL  string `json:"image_url"`
	Position  int    `json:"position"`
	ProductID int    `json:"product_id"`
}

type Address struct {
	ID                  int     `json:"id"`
	Province            string  `json:"province"`
	City                string  `json:"city"`
	District            *string `json:"district"`
	SubDistrict         string  `json:"subdistrict"`
	PostalCode          string  `json:"postal_code"`
	Address             string  `json:"address"`
	ReceivedName        string  `json:"received_name"`
	UserID              string  `json:"user_id"`
	CityVendorID        string  `json:"city_vendor_id"`
	ProvinceVendorID    string  `json:"province_vendor_id"`
	SubDistrictVendorID string  `json:"subdistrict_vendor_id"`
}

type GetWaybillResponse struct {
	Delivered      bool              `json:"delivered"`
	Destination    string            `json:"destination"`
	Resi           string            `json:"resi"`
	ServiceCode    string            `json:"service_code"`
	WaybillDate    string            `json:"waybill_date"`
	ShipperName    string            `json:"shipper_name"`
	ReceiverName   string            `json:"receiver_name"`
	Origin         string            `json:"origin"`
	Status         string            `json:"status"`
	CourierName    string            `json:"courier_name"`
	Manifest       []WaybillManifest `json:"manifest"`
	DeliveryStatus DeliveryStatus    `json:"delivery_status"`
}

type WaybillManifest struct {
	Description string `json:"manifest_description"`
	Date        string `json:"manifest_date"`
	Time        string `json:"manifest_time"`
	City        string `json:"city_name"`
}

type DeliveryStatus struct {
	Status      string `json:"status"`
	PODReceiver string `json:"pod_receiver"`
	PODDate     string `json:"pod_date"`
	PODTime     string `json:"pod_time"`
}

type UpdateOrderShippingNumberRequest struct {
	ShippingNumber string `json:"shipping_number" validate:"required,min=3,max=50,xss_safe"`
}

type UpdateOrderShippingNumberResponse struct {
	ID                  string  `json:"id"`
	OrderDate           string  `json:"order_date"`
	Status              string  `json:"status"`
	ShippingName        string  `json:"shipping_name"`
	ShippingAddress     string  `json:"shipping_address"`
	ShippingPhone       string  `json:"shipping_phone"`
	ShippingNumber      string  `json:"shipping_number"`
	ShippingType        string  `json:"shipping_type"`
	TotalWeight         int     `json:"total_weight"`
	TotalQuantity       int     `json:"total_quantity"`
	TotalShippingCost   string  `json:"total_shipping_cost"`
	TotalProductAmount  string  `json:"total_product_amount"`
	TotalShippingAmount string  `json:"total_shipping_amount"`
	TotalAmount         string  `json:"total_amount"`
	VoucherDiscount     int     `json:"voucher_disc"`
	VoucherID           *string `json:"voucher_id"`
	CostName            string  `json:"cost_name"`
	CostService         string  `json:"cost_service"`
	AddressID           int     `json:"address_id"`
	UserID              string  `json:"user_id"`
	Notes               string  `json:"notes"`
}

type UpdateOrderStatusTransactionRequest struct {
	Status string `json:"status" validate:"required,min=3,max=50,xss_safe"`
}
