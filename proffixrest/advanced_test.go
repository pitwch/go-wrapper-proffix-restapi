package proffixrest

import (
	"context"
	"net/url"
	"testing"
)

// TestClient_Patch tests the PATCH method
func TestClient_Patch(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer pxrest.Logout(ctx)

	// Create a test address first
	data := Adresse{Name: "Patch Test GmbH", Ort: "Bern", PLZ: "3000", Land: Land{LandNr: "CH"}}
	_, headers, statuscode, err := pxrest.Post(ctx, "ADR/Adresse", data)

	if statuscode != 201 {
		t.Errorf("Expected HTTP Status Code 201 for POST. Got '%v'", statuscode)
	}

	if err != nil {
		t.Errorf("Expected no error for POST Request. Got '%v'", err)
	}

	// Get created AdressNr
	adressNr := ConvertLocationToID(headers)

	if adressNr == "" {
		t.Errorf("AdressNr should not be empty")
	}

	// Patch the address (partial update)
	patchData := map[string]interface{}{
		"Name": "Patch Test AG Updated",
	}

	rc, _, patchStatus, patchErr := pxrest.Patch(ctx, "ADR/Adresse/"+adressNr, patchData)

	if patchStatus != 204 {
		t.Errorf("Expected HTTP Status Code 204 for PATCH. Got '%v'", patchStatus)
	}

	if patchErr != nil {
		t.Errorf("Expected no error for PATCH Request. Got '%v'", patchErr)
	}

	if rc != nil {
		rc.Close()
	}

	// Clean up - delete the test address
	_, _, deleteStatus, _ := pxrest.Delete(ctx, "ADR/Adresse/"+adressNr)

	if deleteStatus != 204 {
		t.Errorf("Expected HTTP Status Code 204 for DELETE. Got '%v'", deleteStatus)
	}
}

// TestClient_ServiceLogin tests the ServiceLogin method
func TestClient_ServiceLogin(t *testing.T) {
	ctx := context.Background()

	// First create a normal connection to get a valid session
	pxrest1, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}

	// Get the session ID
	sessionID := pxrest1.GetPxSessionId()

	if sessionID == "" {
		t.Errorf("Expected valid session ID, got empty string")
	}

	// Create a second client and use ServiceLogin
	pxrest2, err := NewClient(
		"https://portal.proffix.net:11011",
		"Gast",
		"16ec7cb001be0525f9af1a96fd5ea26466b2e75ef3e96e881bcb7149cd7598da",
		"DEMODB",
		[]string{"VOL"},
		&Options{
			Key:            "16378f3e3bc8051435694595cbd222219d1ca7f9bddf649b9a0c819a77bb5e50",
			VerifySSL:      false,
			Autologout:     false,
			VolumeLicence:  true,
		},
	)

	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// Use ServiceLogin with the session from first client
	pxrest2.ServiceLogin(ctx, sessionID)

	// Verify the session ID was set
	if pxrest2.GetPxSessionId() != sessionID {
		t.Errorf("Expected session ID '%s', got '%s'", sessionID, pxrest2.GetPxSessionId())
	}

	// Try to use the session (should work)
	_, _, status, err := pxrest2.Get(ctx, "ADR/Adresse", url.Values{})

	if status != 200 {
		t.Errorf("Expected HTTP Status Code 200 with ServiceLogin session. Got '%v'", status)
	}

	// Cleanup
	pxrest1.Logout(ctx)
}

// TestClient_GetPxSessionId tests the GetPxSessionId method
func TestClient_GetPxSessionId(t *testing.T) {
	ctx := context.Background()

	pxrest, err := ConnectTest([]string{})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer pxrest.Logout(ctx)

	sessionID := pxrest.GetPxSessionId()

	if sessionID == "" {
		t.Errorf("Expected non-empty session ID")
	}

	// Verify it's a valid format (should be a string)
	if len(sessionID) < 10 {
		t.Errorf("Session ID seems too short: '%s'", sessionID)
	}
}

// TestClient_UpdatePxSessionId tests that session ID gets updated
func TestClient_UpdatePxSessionId(t *testing.T) {
	ctx := context.Background()

	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer pxrest.Logout(ctx)

	// Get initial session ID
	_ = pxrest.GetPxSessionId()

	// Make a request (this might update the session ID)
	_, _, status, err := pxrest.Get(ctx, "ADR/Adresse", url.Values{})

	if status != 200 {
		t.Errorf("Expected HTTP Status Code 200. Got '%v'", status)
	}

	if err != nil {
		t.Errorf("Expected no error. Got '%v'", err)
	}

	// Get session ID after request
	afterSessionID := pxrest.GetPxSessionId()

	// Session ID should still be valid (either same or updated)
	if afterSessionID == "" {
		t.Errorf("Session ID should not be empty after request")
	}

	// At minimum, we should have a session ID
	if len(afterSessionID) < 10 {
		t.Errorf("Session ID seems invalid: '%s'", afterSessionID)
	}
}

// TestClient_MultipleRequests tests multiple sequential requests
func TestClient_MultipleRequests(t *testing.T) {
	ctx := context.Background()

	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer pxrest.Logout(ctx)

	// Make multiple requests to ensure session handling works
	for i := 0; i < 3; i++ {
		_, _, status, err := pxrest.Get(ctx, "ADR/Adresse", url.Values{})

		if status != 200 {
			t.Errorf("Request %d: Expected HTTP Status Code 200. Got '%v'", i+1, status)
		}

		if err != nil {
			t.Errorf("Request %d: Expected no error. Got '%v'", i+1, err)
		}
	}
}

// TestClient_ConcurrentSafety tests that the client is safe for concurrent use
func TestClient_ConcurrentSafety(t *testing.T) {
	ctx := context.Background()

	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer pxrest.Logout(ctx)

	// Run multiple goroutines accessing the session ID
	done := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		go func(id int) {
			sessionID := pxrest.GetPxSessionId()
			if sessionID == "" {
				t.Errorf("Goroutine %d: Expected non-empty session ID", id)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 5; i++ {
		<-done
	}
}

// TestNewClient_DefaultOptions tests that default options are set correctly
func TestNewClient_DefaultOptions(t *testing.T) {
	client, err := NewClient(
		"https://example.com",
		"user",
		"pass",
		"db",
		[]string{},
		nil,
	)

	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// Check default options
	if client.option.Version != "v4" {
		t.Errorf("Expected default version 'v4', got '%s'", client.option.Version)
	}

	if client.option.APIPrefix != "/pxapi/" {
		t.Errorf("Expected default APIPrefix '/pxapi/', got '%s'", client.option.APIPrefix)
	}

	if client.option.LoginEndpoint != "PRO/Login" {
		t.Errorf("Expected default LoginEndpoint 'PRO/Login', got '%s'", client.option.LoginEndpoint)
	}

	if client.option.Batchsize != 200 {
		t.Errorf("Expected default Batchsize 200, got %d", client.option.Batchsize)
	}

	if !contains(client.option.UserAgent, "go-wrapper-proffix-restapi") {
		t.Errorf("Expected UserAgent to contain 'go-wrapper-proffix-restapi', got '%s'", client.option.UserAgent)
	}
}

// TestNewClient_CustomOptions tests that custom options are respected
func TestNewClient_CustomOptions(t *testing.T) {
	customOptions := &Options{
		Version:       "v3",
		APIPrefix:     "/custom/",
		LoginEndpoint: "Custom/Login",
		UserAgent:     "CustomAgent",
		Batchsize:     500,
		VerifySSL:     true,
		VolumeLicence: true,
	}

	client, err := NewClient(
		"https://example.com",
		"user",
		"pass",
		"db",
		[]string{"ADR"},
		customOptions,
	)

	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// Check custom options are preserved
	if client.option.Version != "v3" {
		t.Errorf("Expected version 'v3', got '%s'", client.option.Version)
	}

	if client.option.Batchsize != 500 {
		t.Errorf("Expected Batchsize 500, got %d", client.option.Batchsize)
	}

	// VolumeLicence should override modules
	if len(client.Module) != 1 || client.Module[0] != "VOL" {
		t.Errorf("Expected Module ['VOL'] with VolumeLicence, got %v", client.Module)
	}
}

// TestNewClient_InvalidURL tests error handling for invalid URL
func TestNewClient_InvalidURL(t *testing.T) {
	_, err := NewClient(
		"://invalid-url",
		"user",
		"pass",
		"db",
		[]string{},
		nil,
	)

	if err == nil {
		t.Errorf("Expected error for invalid URL, got nil")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
