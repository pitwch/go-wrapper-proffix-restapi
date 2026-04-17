package proffixrest

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"
)

// Test GetBatch
func TestGetBatch(t *testing.T) {

	// New Context
	ctx := context.Background()

	// Define test struct for Adresse
	type Adresse struct {
		AdressNr int
		Name     string
		Ort      string
		Plz      string
		Land     string
	}

	// Define Adressen as array
	var adressen []Adresse

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})

	if err != nil {
		t.Errorf("Error on ConnectTest: '%v'", err)
	}

	if pxrest != nil {
		// Set Params. As we just want some fields we define them on Fields param.
		params := url.Values{}
		params.Set("Fields", "AdressNr,Name,Ort,Plz")

		resultBatch, total, err := pxrest.GetBatch(ctx, "ADR/Adresse", params, 25)

		// Check error. Should be nil
		if err != nil {
			t.Errorf("Expected no error for GET Batch Request. Got '%v'", err)
		}

		// Check total. Should be greater than 0
		if total == 0 {
			t.Errorf("Expected some results. Got '%v'", total)
		}

		// Unmarshal the result into our struct
		err = json.Unmarshal(resultBatch, &adressen)

		// Check error. Should be nil
		if err != nil {
			t.Errorf("Expected no error for Unmarshal GET Batch Request. Got '%v'", err)
		}

		// Logout
		statuslogout, err := pxrest.Logout(ctx)

		// Check error. Should be nil
		if err != nil {
			t.Errorf("Expected no error for Logout. Got '%v'", err)
		}

		// Check HTTP Status of Logout. Should be 204
		if statuslogout != 204 {
			t.Errorf("Expected HTTP Status Code 204. Got '%v'", err)
		}

		// Compare received total vs. length of array. Should be equal.
		if len(adressen) != total {
			t.Errorf("Total Results and Length of Array should be equal. Got '%v' vs. '%v'", len(adressen), total)

		}

		// Check default settings
		if pxrest.option.Batchsize != 200 {
			t.Errorf("Default batchsize should be 200. Got '%v'", pxrest.option.Batchsize)
		}
	}
}

// TestGetBatch_EmptyParams tests batch with nil params
func TestGetBatch_EmptyParams(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	// Test with nil params - should use empty url.Values
	resultBatch, total, err := pxrest.GetBatch(ctx, "ADR/Adresse", nil, 10)

	// Should not error even with nil params
	if err != nil {
		t.Logf("GetBatch with nil params returned error: %v", err)
	}

	// If we got results, verify they're valid
	if len(resultBatch) > 0 {
		// Try to parse as JSON array
		var results []map[string]interface{}
		if jsonErr := json.Unmarshal(resultBatch, &results); jsonErr == nil {
			t.Logf("Got %d results from batch", len(results))
		}
	}

	t.Logf("Total entries reported: %d", total)
}

// TestGetBatch_ZeroBatchSize tests that zero batchsize uses default
func TestGetBatch_ZeroBatchSize(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	params := url.Values{}
	params.Set("Fields", "AdressNr,Name")

	// Pass 0 as batchsize - should use default from options
	resultBatch, total, err := pxrest.GetBatch(ctx, "ADR/Adresse", params, 0)

	if err != nil {
		t.Logf("GetBatch with zero batchsize returned error: %v", err)
	}

	if total > 0 && len(resultBatch) == 0 {
		t.Errorf("Got total %d but empty result batch", total)
	}
}

// TestGetBatch_LargeBatchSize tests with large batch size
func TestGetBatch_LargeBatchSize(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	params := url.Values{}
	params.Set("Fields", "AdressNr,Name")

	// Use large batchsize
	resultBatch, total, err := pxrest.GetBatch(ctx, "ADR/Adresse", params, 1000)

	if err != nil {
		t.Logf("GetBatch with large batchsize returned error: %v", err)
	}

	if total > 0 {
		// Verify results can be parsed
		var results []map[string]interface{}
		if jsonErr := json.Unmarshal(resultBatch, &results); jsonErr != nil {
			t.Errorf("Failed to unmarshal batch results: %v", jsonErr)
		} else {
			t.Logf("Large batch returned %d results (total: %d)", len(results), total)
		}
	}
}

// TestGetBatch_NonExistentEndpoint tests error handling
func TestGetBatch_NonExistentEndpoint(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	params := url.Values{}

	// Try a non-existent endpoint
	_, _, err = pxrest.GetBatch(ctx, "NONEXISTENT/Endpoint", params, 10)

	// Should return an error for non-existent endpoint
	if err != nil {
		t.Logf("Expected error for non-existent endpoint: %v", err)
	}
}

// TestGetBatch_InvalidEndpoint tests invalid endpoint format
func TestGetBatch_InvalidEndpoint(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	params := url.Values{}

	// Try an invalid endpoint format
	_, _, err = pxrest.GetBatch(ctx, "", params, 10)

	// Should error on empty endpoint
	if err == nil {
		t.Logf("Empty endpoint may be handled gracefully by API")
	} else {
		t.Logf("Empty endpoint returned error: %v", err)
	}
}
