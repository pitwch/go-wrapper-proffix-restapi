package proffixrest

import (
	"net/http"
	"path"
	"time"
)

//Time Constant for parsing Time Format of PROFFIX REST-API
//Used for Function ConvertPXTime
const PXTime = "2006-01-02 15:04:05"

//Convert Timestamps of PROFFIX REST-API to time.Time
//Accepts time as string from PROFFIX REST-API
//Returns Golang time.time
func ConvertPXTimeToTime(pxtime string) (t time.Time) {
	t, _ = time.Parse(PXTime, pxtime)
	return t
}

//Convert time.Time to PROFFIX REST-API Timestamp
//Accepts time.Time
//Returns PROFFX REST-API Timestamp as string
func ConvertTimeToPXTime(t time.Time) (pxtime string) {
	return t.Format(PXTime)
}

//Converts Header of PROFFIX REST-API to ID
//Used for Get the ID of new created element through POST Method
//Accepts Header as http.Header
//Returns ID as string
func ConvertLocationToID(header http.Header) (id string) {
	return path.Base(header.Get("Location"))
}
