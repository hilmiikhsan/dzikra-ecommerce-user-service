package dto

type GetListSubDistrictResponse struct {
	ID              string `json:"id"`
	CityName        string `json:"city_name"`
	CityType        string `json:"city_type"`
	ProvinceName    string `json:"province_name"`
	SubDistrictName string `json:"subdistrict_name"`
}
