package entity

import (
	"database/sql"

	"github.com/google/uuid"
)

type Address struct {
	ID                  int            `db:"id"`
	Province            string         `db:"province"`
	City                string         `db:"city"`
	District            sql.NullString `db:"district"`
	SubDistrict         string         `db:"subdistrict"`
	PostalCode          string         `db:"postal_code"`
	Address             string         `db:"address"`
	ReceivedName        string         `db:"received_name"`
	UserID              uuid.UUID      `db:"user_id"`
	CityVendorID        string         `db:"city_vendor_id"`
	ProvinceVendorID    string         `db:"province_vendor_id"`
	SubDistrictVendorID string         `db:"subdistrict_vendor_id"`
}
