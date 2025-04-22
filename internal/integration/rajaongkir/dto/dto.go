package dto

type RajaOngkirProvincePayload struct {
	Rajaongkir struct {
		Results []ProvinceResult `json:"results"`
	} `json:"rajaongkir"`
}

type ProvinceResult struct {
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
}

type RajaOngkirCityPayload struct {
	Rajaongkir struct {
		Results []CityResult `json:"results"`
	} `json:"rajaongkir"`
}

type CityResult struct {
	CityID       string `json:"city_id"`
	CityName     string `json:"city_name"`
	Type         string `json:"type"`
	ProvinceName string `json:"province"`
	PostalCode   string `json:"postal_code"`
}
