package helpers

import (
	"log"
	"os"
	"time"
)

func GetAppLocation() *time.Location {
	tz := os.Getenv("APP_TIMEZONE")
	if tz == "" {
		tz = "UTC"
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Printf("failed load timezone %s, fallback to UTC", tz)
		return time.UTC
	}

	return loc
}
