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

type RajaOngkirSubDistrictPayload struct {
	Rajaongkir struct {
		Results []SubDistrictResult `json:"results"`
	} `json:"rajaongkir"`
}

type SubDistrictResult struct {
	SubDistrictID   string `json:"subdistrict_id"`
	CityName        string `json:"city_name"`
	CityType        string `json:"type"`
	ProvinceName    string `json:"province"`
	SubDistrictName string `json:"subdistrict_name"`
}

type RajaOngkirCostPayload struct {
	Rajaongkir struct {
		Results []CostResult `json:"results"`
	} `json:"rajaongkir"`
}

type CostResult struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Cost []struct {
		Service     string `json:"service"`
		Description string `json:"description"`
		Cost        []struct {
			Value int    `json:"value"`
			Etd   string `json:"etd"`
			Note  string `json:"note"`
		} `json:"cost"`
	} `json:"costs"`
}
