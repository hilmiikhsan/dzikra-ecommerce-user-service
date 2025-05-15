package dto

type CalculateShippingCostRequest struct {
	Courier   string `json:"courier" validate:"required,min=3,max=20,xss_safe"`
	AddressID int    `json:"address_id" validate:"required,numeric,gt=0"`
}

type CalculateShippingCostResponse struct {
	Code  string               `json:"code"`
	Name  string               `json:"name"`
	Costs []DetailShippingCost `json:"costs"`
}

type DetailShippingCost struct {
	Service     string       `json:"service"`
	Description string       `json:"description"`
	Cost        []DetailCost `json:"cost"`
}

type DetailCost struct {
	Value int    `json:"value"`
	Etd   string `json:"etd"`
	Note  string `json:"note"`
}
