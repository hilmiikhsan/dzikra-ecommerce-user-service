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

type RajaOngkirWaybillPayload struct {
	Rajaongkir struct {
		Result struct {
			Summary struct {
				WaybillNumber string `json:"waybill_number"`
				ServiceCode   string `json:"service_code"`
				WaybillDate   string `json:"waybill_date"`
				ShipperName   string `json:"shipper_name"`
				ReceiverName  string `json:"receiver_name"`
				Origin        string `json:"origin"`
				Destination   string `json:"destination"`
				Status        string `json:"status"`
				CourierName   string `json:"courier_name"`
			} `json:"summary"`
			Manifest []struct {
				ManifestDescription string `json:"manifest_description"`
				ManifestDate        string `json:"manifest_date"`
				ManifestTime        string `json:"manifest_time"`
				CityName            string `json:"city_name"`
			} `json:"manifest"`
			DeliveryStatus struct {
				Status      string `json:"status"`
				PODReceiver string `json:"pod_receiver"`
				PODDate     string `json:"pod_date"`
				PODTime     string `json:"pod_time"`
			} `json:"delivery_status"`
		} `json:"result"`
	} `json:"rajaongkir"`
}

// DTO yang akan dikirim ke handler/API layer
type GetWaybillResponse struct {
	Summary        WaybillSummary        `json:"summary"`
	Manifest       []WaybillManifest     `json:"manifest"`
	DeliveryStatus WaybillDeliveryStatus `json:"delivery_status"`
}

type WaybillSummary struct {
	Resi         string `json:"resi"`
	ServiceCode  string `json:"service_code"`
	WaybillDate  string `json:"waybill_date"`
	ShipperName  string `json:"shipper_name"`
	ReceiverName string `json:"receiver_name"`
	Origin       string `json:"origin"`
	Destination  string `json:"destination"`
	Status       string `json:"status"`
	CourierName  string `json:"courier_name"`
}

type WaybillManifest struct {
	Description string `json:"manifest_description"`
	Date        string `json:"manifest_date"`
	Time        string `json:"manifest_time"`
	CityName    string `json:"city_name"`
}

type WaybillDeliveryStatus struct {
	Status      string `json:"status"`
	PODReceiver string `json:"pod_receiver"`
	PODDate     string `json:"pod_date"`
	PODTime     string `json:"pod_time"`
}
