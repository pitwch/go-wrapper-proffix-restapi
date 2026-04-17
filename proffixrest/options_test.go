package proffixrest

import (
	"net/http"
	"testing"
	"time"
)

func TestOptions_DefaultValues(t *testing.T) {
	// Create client with nil options to test defaults
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

	// Check default Version
	if client.option.Version != "v4" {
		t.Errorf("Expected default Version 'v4', got '%s'", client.option.Version)
	}

	// Check default APIPrefix
	if client.option.APIPrefix != "/pxapi/" {
		t.Errorf("Expected default APIPrefix '/pxapi/', got '%s'", client.option.APIPrefix)
	}

	// Check default LoginEndpoint
	if client.option.LoginEndpoint != "PRO/Login" {
		t.Errorf("Expected default LoginEndpoint 'PRO/Login', got '%s'", client.option.LoginEndpoint)
	}

	// Check default Batchsize
	if client.option.Batchsize != 200 {
		t.Errorf("Expected default Batchsize 200, got %d", client.option.Batchsize)
	}

	// Check UserAgent contains wrapper name
	if client.option.UserAgent == "" {
		t.Errorf("Expected non-empty UserAgent")
	}
}

func TestOptions_CustomValues(t *testing.T) {
	customTimeout := 30 * time.Second
	customHTTPClient := &http.Client{Timeout: customTimeout}

	customOptions := &Options{
		Key:           "custom-key",
		Version:       "v3",
		APIPrefix:     "/custom/api/",
		LoginEndpoint: "Auth/Login",
		UserAgent:     "CustomUserAgent/1.0",
		Timeout:       customTimeout,
		Batchsize:     500,
		Log:           true,
		Autologout:    true,
		VerifySSL:     true,
		VolumeLicence: false,
		HTTPClient:    customHTTPClient,
	}

	client, err := NewClient(
		"https://example.com",
		"user",
		"pass",
		"db",
		[]string{"ADR", "LAG"},
		customOptions,
	)

	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// Verify custom values are preserved
	if client.option.Key != "custom-key" {
		t.Errorf("Expected Key 'custom-key', got '%s'", client.option.Key)
	}

	if client.option.Version != "v3" {
		t.Errorf("Expected Version 'v3', got '%s'", client.option.Version)
	}

	if client.option.APIPrefix != "/custom/api/" {
		t.Errorf("Expected APIPrefix '/custom/api/', got '%s'", client.option.APIPrefix)
	}

	if client.option.LoginEndpoint != "Auth/Login" {
		t.Errorf("Expected LoginEndpoint 'Auth/Login', got '%s'", client.option.LoginEndpoint)
	}

	if client.option.UserAgent != "CustomUserAgent/1.0" {
		t.Errorf("Expected UserAgent 'CustomUserAgent/1.0', got '%s'", client.option.UserAgent)
	}

	if client.option.Timeout != customTimeout {
		t.Errorf("Expected Timeout %v, got %v", customTimeout, client.option.Timeout)
	}

	if client.option.Batchsize != 500 {
		t.Errorf("Expected Batchsize 500, got %d", client.option.Batchsize)
	}

	if !client.option.Log {
		t.Errorf("Expected Log to be true")
	}

	if !client.option.Autologout {
		t.Errorf("Expected Autologout to be true")
	}

	if !client.option.VerifySSL {
		t.Errorf("Expected VerifySSL to be true")
	}

	// Verify HTTPClient is used
	if client.option.HTTPClient != customHTTPClient {
		t.Errorf("Expected custom HTTPClient to be preserved")
	}

	// Modules should NOT be overridden since VolumeLicence is false
	if len(client.Module) != 2 || client.Module[0] != "ADR" {
		t.Errorf("Expected Modules ['ADR', 'LAG'], got %v", client.Module)
	}
}

func TestOptions_VolumeLicenceOverridesModules(t *testing.T) {
	customOptions := &Options{
		VolumeLicence: true,
	}

	client, err := NewClient(
		"https://example.com",
		"user",
		"pass",
		"db",
		[]string{"ADR", "LAG", "FIB"},
		customOptions,
	)

	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// VolumeLicence should override modules to just VOL
	if len(client.Module) != 1 || client.Module[0] != "VOL" {
		t.Errorf("Expected Module ['VOL'] with VolumeLicence, got %v", client.Module)
	}
}

func TestOptions_DefaultHTTPClient(t *testing.T) {
	// Create client without custom HTTPClient
	client, err := NewClient(
		"https://example.com",
		"user",
		"pass",
		"db",
		[]string{},
		&Options{VerifySSL: false},
	)

	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// Should have created a default HTTP client
	if client.option.HTTPClient != nil {
		t.Errorf("Expected nil HTTPClient when not provided")
	}

	// The client's internal http.Client should be set
	if client.client == nil {
		t.Errorf("Expected internal http.Client to be created")
	}
}

func TestOptions_EmptyOptions(t *testing.T) {
	// Test with empty Options struct
	emptyOptions := &Options{}

	client, err := NewClient(
		"https://example.com",
		"user",
		"pass",
		"db",
		[]string{},
		emptyOptions,
	)

	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// Should still apply defaults for empty fields
	if client.option.Version != "v4" {
		t.Errorf("Expected default Version 'v4', got '%s'", client.option.Version)
	}

	if client.option.APIPrefix != "/pxapi/" {
		t.Errorf("Expected default APIPrefix '/pxapi/', got '%s'", client.option.APIPrefix)
	}

	if client.option.Batchsize != 200 {
		t.Errorf("Expected default Batchsize 200, got %d", client.option.Batchsize)
	}
}
