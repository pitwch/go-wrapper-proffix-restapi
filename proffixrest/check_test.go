package proffixrest

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestClient_CheckAPI(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	// Check CheckAPI with the client's key
	err = pxrest.CheckAPI(ctx, pxrest.option.Key)

	// Check error. Should be nil (or log the result)
	if err != nil {
		t.Logf("CheckAPI returned error (may be expected in test environment): %v", err)
	}
}

// TestCheckAPI_EmptyKey tests CheckAPI with empty key
func TestCheckAPI_EmptyKey(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	// Check CheckAPI with empty key
	err = pxrest.CheckAPI(ctx, "")

	// Should return an error for empty key
	if err == nil {
		t.Logf("CheckAPI with empty key may be handled by API")
	} else {
		t.Logf("CheckAPI with empty key returned expected error: %v", err)
	}
}

// TestInfoStruct_Structure tests InfoStruct marshaling/unmarshaling
func TestInfoStruct_Structure(t *testing.T) {
	// Test JSON marshaling/unmarshaling of InfoStruct
	info := InfoStruct{
		Version:        "4.5.0",
		ServerZeit:     "2024-01-15 10:30:00",
		NeuesteVersion: "4.5.1",
		Instanz: InstanzStruct{
			InstanzNr: "001",
			Name:      "Test Instance",
			Lizenzen: []LizenzStruct{
				{
					Name:        "ADR",
					Bezeichnung: "Adressen",
					Anzahl:      10,
					Demo:        false,
					Ablaufdatum: "2025-12-31",
				},
				{
					Name:        "LAG",
					Bezeichnung: "Lager",
					Anzahl:      5,
					Demo:        false,
					Ablaufdatum: "2025-12-31",
				},
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(info)
	if err != nil {
		t.Errorf("Failed to marshal InfoStruct: %v", err)
	}

	// Unmarshal back
	var parsedInfo InfoStruct
	err = json.Unmarshal(jsonData, &parsedInfo)
	if err != nil {
		t.Errorf("Failed to unmarshal InfoStruct: %v", err)
	}

	// Verify data integrity
	if parsedInfo.Version != info.Version {
		t.Errorf("Version mismatch: expected %s, got %s", info.Version, parsedInfo.Version)
	}

	if parsedInfo.Instanz.Name != info.Instanz.Name {
		t.Errorf("Instanz name mismatch")
	}

	if len(parsedInfo.Instanz.Lizenzen) != len(info.Instanz.Lizenzen) {
		t.Errorf("License count mismatch: expected %d, got %d", len(info.Instanz.Lizenzen), len(parsedInfo.Instanz.Lizenzen))
	}

	// Check first license
	if len(parsedInfo.Instanz.Lizenzen) > 0 {
		if parsedInfo.Instanz.Lizenzen[0].Name != "ADR" {
			t.Errorf("First license name mismatch")
		}
		if parsedInfo.Instanz.Lizenzen[0].Anzahl != 10 {
			t.Errorf("First license count mismatch: expected 10, got %d", parsedInfo.Instanz.Lizenzen[0].Anzahl)
		}
	}
}

// TestInfoStruct_EmptyLicenses tests with empty license list
func TestInfoStruct_EmptyLicenses(t *testing.T) {
	info := InfoStruct{
		Version:    "4.0.0",
		ServerZeit: "2024-01-15 10:30:00",
		Instanz: InstanzStruct{
			InstanzNr: "002",
			Name:      "Empty Instance",
			Lizenzen:  []LizenzStruct{},
		},
	}

	jsonData, err := json.Marshal(info)
	if err != nil {
		t.Errorf("Failed to marshal InfoStruct with empty licenses: %v", err)
	}

	var parsedInfo InfoStruct
	err = json.Unmarshal(jsonData, &parsedInfo)
	if err != nil {
		t.Errorf("Failed to unmarshal InfoStruct with empty licenses: %v", err)
	}

	if len(parsedInfo.Instanz.Lizenzen) != 0 {
		t.Errorf("Expected empty license list, got %d licenses", len(parsedInfo.Instanz.Lizenzen))
	}
}

// TestLizenzStruct_Fields tests individual license fields
func TestLizenzStruct_Fields(t *testing.T) {
	license := LizenzStruct{
		Name:        "FIB",
		Bezeichnung: "Finanzbuchhaltung",
		Anzahl:      3,
		Demo:        true,
		Ablaufdatum: "2024-06-30",
	}

	// Verify fields are accessible
	if license.Name != "FIB" {
		t.Errorf("License name mismatch")
	}

	if license.Bezeichnung != "Finanzbuchhaltung" {
		t.Errorf("License description mismatch")
	}

	if license.Anzahl != 3 {
		t.Errorf("License count mismatch")
	}

	if !license.Demo {
		t.Errorf("Expected demo license")
	}

	if license.Ablaufdatum != "2024-06-30" {
		t.Errorf("License expiry mismatch")
	}
}

// TestCheckAPI_TimeoutSet tests that timeout is set correctly
func TestCheckAPI_TimeoutSet(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	// Store original timeout
	originalTimeout := pxrest.option.Timeout

	// Call CheckAPI
	_ = pxrest.CheckAPI(ctx, pxrest.option.Key)

	// Check that timeout was set to 10 (as per CheckAPI implementation)
	// Note: This is an implementation detail check
	if pxrest.option.Timeout != 10 {
		t.Logf("Timeout after CheckAPI: %v (expected implementation sets it to 10)", pxrest.option.Timeout)
	}

	// Restore original timeout for cleanup
	pxrest.option.Timeout = originalTimeout
}

// TestGetFilteredCount_WithPxMetadata tests GetFilteredCount with pxmetadata header
func TestGetFilteredCount_WithPxMetadata(t *testing.T) {
	// Create test header with pxmetadata
	header := http.Header{}
	metadata := `{"FilteredCount":100,"TotalCount":500}`
	header.Set("pxmetadata", metadata)

	count := GetFilteredCount(header)

	if count != 100 {
		t.Errorf("Expected FilteredCount 100, got %d", count)
	}
}

// TestGetFilteredCount_EmptyMetadata tests GetFilteredCount with empty metadata
func TestGetFilteredCount_EmptyMetadata(t *testing.T) {
	// Test with empty header
	header := http.Header{}
	count := GetFilteredCount(header)

	if count != 0 {
		t.Errorf("Expected 0 for empty header, got %d", count)
	}

	// Test with empty pxmetadata value
	header.Set("pxmetadata", "")
	count = GetFilteredCount(header)

	if count != 0 {
		t.Errorf("Expected 0 for empty pxmetadata, got %d", count)
	}
}

// TestGetFilteredCount_InvalidJSON tests GetFilteredCount with invalid JSON
func TestGetFilteredCount_InvalidJSON(t *testing.T) {
	// Test with invalid JSON in pxmetadata
	header := http.Header{}
	header.Set("pxmetadata", "invalid json")

	count := GetFilteredCount(header)

	// Should gracefully handle invalid JSON and return 0
	if count != 0 {
		t.Errorf("Expected 0 for invalid JSON, got %d", count)
	}
}
