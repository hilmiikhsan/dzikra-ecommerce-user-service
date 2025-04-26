package dto

type GetListCityResponse struct {
	ID           string `json:"id"`
	City         string `json:"city"`
	Type         string `json:"type"`
	ProvinceName string `json:"province_name"`
	PostalCode   string `json:"postal_code"`
}
