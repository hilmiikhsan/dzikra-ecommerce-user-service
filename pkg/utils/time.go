package utils

import "time"

func FormatToWIB(t time.Time) string {
	// WIB (UTC+7)
	wibZone := time.FixedZone("WIB", 7*3600)
	return t.In(wibZone).Format("2006-01-02T15:04:05-07:00")
}
