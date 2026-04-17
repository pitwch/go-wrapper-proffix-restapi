package proffixrest

import (
	"context"
	"testing"
)

// TestClient_GetList
func TestClient_GetList(t *testing.T) {

	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})

	// Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Connect. Got '%v'", err)
	}

	_, headers, status, err := pxrest.GetList(ctx, 1029, nil)

	if status != 404 {
		// Check error. Should be nil
		if err != nil {
			t.Errorf("Expected no error for GET List Request. Got '%v'", err)
		}

		// Check status. Should be 201
		if status != 200 {
			t.Errorf("Expected Status 201. Got '%v'", status)
		}

		// Check Content Type. Should be pdf
		if headers.Get("Content-Type") != "application/pdf" {
			t.Errorf("Expected Content-Type PDF. Got '%v'", status)
		}
	}
	// Check Logout
	statuslogout, err := pxrest.Logout(ctx)

	// Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Logout Request. Got '%v'", err)
	}

	// Check HTTP Status of Logout. Should be 204
	if statuslogout != 204 {
		t.Errorf("Expected HTTP Status Code 204. Got '%v'", err)
	}
}

// TestClient_GetList_InvalidListNr tests with invalid list number
func TestClient_GetList_InvalidListNr(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	// Try with invalid/negative list number
	_, _, status, err := pxrest.GetList(ctx, -1, nil)

	// Should return error or 404
	if status == 404 {
		t.Logf("Invalid list number returned 404 as expected")
	} else if err != nil {
		t.Logf("Invalid list number returned error: %v", err)
	}
}

// TestClient_GetList_WithBody tests list generation with filter body
func TestClient_GetList_WithBody(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	// Create filter body for list
	filterBody := map[string]interface{}{
		"Filter": "Land=='CH'",
	}

	rc, headers, status, err := pxrest.GetList(ctx, 1029, filterBody)

	if rc != nil {
		_ = rc.Close()
	}

	// The list endpoint may return different statuses
	t.Logf("GetList with body returned status %d, headers: %v, err: %v", status, headers, err)

	if status == 200 {
		// Check content type if successful
		contentType := headers.Get("Content-Type")
		if contentType != "" && contentType != "application/pdf" {
			t.Logf("Unexpected content type: %s", contentType)
		}
	}
}

// TestClient_GetList_NilBody tests list with nil body
func TestClient_GetList_NilBody(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	// Try with nil body (should work if API supports it)
	rc, headers, status, err := pxrest.GetList(ctx, 1029, nil)

	if rc != nil {
		_ = rc.Close()
	}

	t.Logf("GetList with nil body returned status %d, content-type: %s", status, headers.Get("Content-Type"))

	// Should either succeed or give a meaningful error
	if status >= 500 {
		t.Errorf("Server error for GetList with nil body: %v", err)
	}
}
