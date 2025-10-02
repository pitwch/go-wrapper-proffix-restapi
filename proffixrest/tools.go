package proffixrest

import (
	"net/http"
	"path"
	"time"
)

// PXTime is the layout used to parse PROFFIX REST-API timestamps.
const PXTime = "2006-01-02 15:04:05"

// ConvertPXTimeToTime converts a PROFFIX REST-API timestamp string to time.Time.
func ConvertPXTimeToTime(pxtime string) (t time.Time) {
	t, _ = time.Parse(PXTime, pxtime) // Ignore error, return zero time on failure
	return t
}

// ConvertTimeToPXTime converts a time.Time to the PROFFIX REST-API timestamp format.
func ConvertTimeToPXTime(t time.Time) (pxtime string) {
	return t.Format(PXTime)
}

// ConvertLocationToID extracts the trailing resource identifier from a Location header.
func ConvertLocationToID(header http.Header) (id string) {
	return path.Base(header.Get("Location"))
}
