package proffixrest

import (
	"net/http"
	"testing"
	"time"
)

func TestConvertLocationToID(t *testing.T) {
	// Create test header with Location
	header := http.Header{}
	header.Set("Location", "https://example.com/api/v1/ADR/Adresse/12345")

	// Test conversion
	id := ConvertLocationToID(header)

	// Check if ID is correct
	if id != "12345" {
		t.Errorf("Expected ID '12345', got '%s'", id)
	}

	// Test with empty header
	emptyHeader := http.Header{}
	emptyID := ConvertLocationToID(emptyHeader)
	if emptyID != "." {
		t.Errorf("Expected '.' for empty header, got '%s'", emptyID)
	}
}

func TestTimeConversion(t *testing.T) {

	// Set testtime as time.Now()
	testtime := time.Now()

	// Convert testtime to pxtime
	pxtime := ConvertTimeToPXTime(testtime)

	// Convert back: pxtime into time
	gotime := ConvertPXTimeToTime(pxtime)

	// Check if still same time
	if testtime.Day() != gotime.Day() {
		t.Errorf("Testtime %v and Gotime %v should be same", testtime, gotime)
	}

}
