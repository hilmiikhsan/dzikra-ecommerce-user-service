package entity

import (
	"database/sql"

	"github.com/google/uuid"
)

type Address struct {
	ID                  int            `json:"id"`
	Province            string         `json:"province"`
	City                string         `json:"city"`
	District            sql.NullString `json:"district"`
	SubDistrict         string         `json:"subdistrict"`
	PostalCode          string         `json:"postal_code"`
	Address             string         `json:"address"`
	ReceivedName        string         `json:"received_name"`
	UserID              uuid.UUID      `json:"user_id"`
	CityVendorID        string         `json:"city_vendor_id"`
	ProvinceVendorID    string         `json:"province_vendor_id"`
	SubDistrictVendorID string         `json:"subdistrict_vendor_id"`
}
