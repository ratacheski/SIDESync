package utils

import (
	"time"
)

func FormataDataHoraWebService(data time.Time) string {
	format := data.Format("20060102150405")
	return format
}
