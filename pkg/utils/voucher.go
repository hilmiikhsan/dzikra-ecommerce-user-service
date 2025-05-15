package utils

import (
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/constants"
)

func IsVoucherExpired(endAt time.Time) bool {
	loc, err := time.LoadLocation(constants.TimeLocationAsiaJakarta)
	if err != nil {
		loc = time.UTC
	}

	now := time.Now().In(loc)
	end := endAt.In(loc)

	return now.After(end)
}
