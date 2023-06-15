package helpers

import (
	"time"
)

func UnixToTime(unixTime int64) string {
	// Create a time object from Unix timestamp
	t := time.Unix(unixTime, 0)

	// Format the time object to desired layout
	layout := "2006-01-02 15:04:05" // Use any layout you prefer
	humanReadableTime := t.Format(layout)

	return humanReadableTime
}
