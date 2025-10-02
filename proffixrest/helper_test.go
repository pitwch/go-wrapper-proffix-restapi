package proffixrest

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func Test_errorFormatterPx(t *testing.T) {
	ctx := context.Background()
	pxrest, _ := ConnectTest([]string{})

	_, _, status, err := pxrest.Get(ctx, "ADR/Adresse/12345678", nil)
	_, _ = pxrest.Logout(ctx)

	if status != 404 {
		t.Errorf("Statuscode should be 404: Statuscode %v", status)
	}

	if !err.(*PxError).isNotFound() {
		t.Errorf("Error should be of Type 'NOT_FOUND'")
	}
}

func TestGetFilteredCount(t *testing.T) {
	// Create test header with pxmetadata
	header := http.Header{}
	metadata := `{"FilteredCount":42}`
	header.Set("pxmetadata", metadata)

	count := GetFilteredCount(header)

	if count != 42 {
		t.Errorf("Expected FilteredCount 42, got %v", count)
	}

	// Test with empty header
	emptyHeader := http.Header{}
	emptyCount := GetFilteredCount(emptyHeader)
	if emptyCount != 0 {
		t.Errorf("Expected FilteredCount 0 for empty header, got %v", emptyCount)
	}
}

func TestReaderToString(t *testing.T) {
	testString := "Hello, PROFFIX REST API!"
	reader := strings.NewReader(testString)

	result, err := ReaderToString(reader)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != testString {
		t.Errorf("Expected '%s', got '%s'", testString, result)
	}

	// Test with nil reader
	nilResult, nilErr := ReaderToString(nil)
	if nilErr != nil {
		t.Errorf("Expected no error for nil reader, got %v", nilErr)
	}
	if nilResult != "" {
		t.Errorf("Expected empty string for nil reader, got '%s'", nilResult)
	}
}

func TestReaderToByte(t *testing.T) {
	testData := []byte("Test data for PROFFIX")
	reader := bytes.NewReader(testData)

	result, err := ReaderToByte(reader)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !bytes.Equal(result, testData) {
		t.Errorf("Expected %v, got %v", testData, result)
	}

	// Test with nil reader
	nilResult, nilErr := ReaderToByte(nil)
	if nilErr != nil {
		t.Errorf("Expected no error for nil reader, got %v", nilErr)
	}
	if nilResult != nil {
		t.Errorf("Expected nil for nil reader, got %v", nilResult)
	}
}

func TestGetMaps(t *testing.T) {
	testJSON := `[{"Name":"Test1","Value":1},{"Name":"Test2","Value":2}]`
	reader := strings.NewReader(testJSON)

	result, err := GetMaps(reader)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 items, got %v", len(result))
	}

	if result[0]["Name"] != "Test1" {
		t.Errorf("Expected Name 'Test1', got %v", result[0]["Name"])
	}
}

func TestGetMap(t *testing.T) {
	testJSON := `{"Name":"TestItem","Value":42}`
	reader := strings.NewReader(testJSON)

	result, err := GetMap(reader)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result["Name"] != "TestItem" {
		t.Errorf("Expected Name 'TestItem', got %v", result["Name"])
	}

	if result["Value"].(float64) != 42 {
		t.Errorf("Expected Value 42, got %v", result["Value"])
	}
}

func TestWriteFile(t *testing.T) {
	// Create temp file path
	tempFile := "test_write_file.txt"
	defer func() { _ = os.Remove(tempFile) }()

	// Create test data
	testData := "Test file content for PROFFIX REST API"
	reader := io.NopCloser(strings.NewReader(testData))

	// Write file
	err := WriteFile(tempFile, reader)
	if err != nil {
		t.Errorf("Expected no error writing file, got %v", err)
	}

	// Read back and verify
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Errorf("Expected no error reading file, got %v", err)
	}

	if string(content) != testData {
		t.Errorf("Expected '%s', got '%s'", testData, string(content))
	}
}

func TestGetUsedLicences(t *testing.T) {
	testJSON := `{
		"Instanz": {
			"Lizenzen": [
				{"Name":"ADR","AnzahlInVerwendung":5},
				{"Name":"FIB","AnzahlInVerwendung":3},
				{"Name":"LAG","AnzahlInVerwendung":2}
			]
		}
	}`
	reader := strings.NewReader(testJSON)

	total, err := GetUsedLicences(reader)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if total != 10 {
		t.Errorf("Expected total 10, got %v", total)
	}
}

func TestGetFileTokens(t *testing.T) {
	testJSON := `[
		{"ArtikelNr":"123","DateiNr":"456"},
		{"ArtikelNr":"789","DateiNr":"012"}
	]`
	reader := strings.NewReader(testJSON)

	files, err := GetFileTokens(reader, "ArtikelNr", "DateiNr")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 file tokens, got %v", len(files))
	}

	// Check first item
	if files[0]["123"] == nil {
		t.Errorf("Expected key '123' in first item")
	}
}
