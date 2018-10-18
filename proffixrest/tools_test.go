package proffixrest

import (
	"testing"
	"time"
)

func TestConvertLocationToID(t *testing.T) {

}

func TestTimeConversion(t *testing.T) {

	//Set testtime as time.Now()
	testtime := time.Now()

	//Convert testtime to pxtime
	pxtime := ConvertTimeToPXTime(testtime)

	//Convert back: pxtime into time
	gotime := ConvertPXTimeToTime(pxtime)

	//Check if still same time
	if testtime.Day() != gotime.Day() {
		t.Errorf("Testtime %v and Gotime %v should be same", testtime, gotime)
	}

}
