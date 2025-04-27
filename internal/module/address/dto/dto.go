package dto

type CreateOrUpdateAddressRequest struct {
	ReceivedName  string `json:"received_name" validate:"required,max=100,min=2,xss_safe"`
	Province      string `json:"province" validate:"required,max=100,min=2,xss_safe"`
	ProvinceID    string `json:"province_id" validate:"required,max=100,min=1,xss_safe"`
	City          string `json:"city" validate:"required,max=100,min=2,xss_safe"`
	CityID        string `json:"city_id" validate:"required,max=100,min=1,xss_safe"`
	SubDistrict   string `json:"subdistrict" validate:"required,max=100,min=2,xss_safe"`
	SubDistrictID string `json:"subdistrict_id" validate:"required,max=100,min=1,xss_safe"`
	Address       string `json:"address" validate:"required,max=255,min=2,xss_safe"`
	PostalCode    string `json:"postal_code" validate:"required,max=10,min=2,xss_safe"`
	UserID        string `json:"user_id" validate:"required,uuid"`
}

type CreateOrUpdateAddressResponse struct {
	ID                  int     `json:"id"`
	Province            string  `json:"province"`
	ProvinceVendorID    string  `json:"province_vendor_id"`
	City                string  `json:"city"`
	CityVendorID        string  `json:"city_vendor_id"`
	District            *string `json:"district"`
	SubDistrict         string  `json:"subdistrict"`
	SubDistrictVendorID string  `json:"subdistrict_vendor_id"`
	Address             string  `json:"address"`
	PostalCode          string  `json:"postal_code"`
	ReceivedName        string  `json:"received_name"`
	UserID              string  `json:"user_id"`
}

type GetListAddressResponse struct {
	ID                  int    `json:"id"`
	Province            string `json:"province"`
	ProvinceVendorID    string `json:"province_vendor_id"`
	City                string `json:"city"`
	CityVendorID        string `json:"city_vendor_id"`
	SubDistrict         string `json:"subdistrict"`
	SubDistrictVendorID string `json:"subdistrict_vendor_id"`
	Address             string `json:"address"`
	PostalCode          string `json:"postal_code"`
	ReceivedName        string `json:"received_name"`
}
