package proffixrest

import "time"

//Define Time Constant for parsing Time Format
const PXTime = "2006-02-01 15:04:05"

//Convert Timestamps of PROFFIX REST-API to time.Time
func ConvertPXTime(pxtime string) (t time.Time) {
	t, _ = time.Parse(PXTime, pxtime)
	return t
}
