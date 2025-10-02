package proffixrest

import (
	"strings"
	"testing"
)

func TestPxError_isInvalidFields(t *testing.T) {
	// Test with INVALID_FIELDS type
	err := &PxError{
		Type:    "INVALID_FIELDS",
		Message: "Invalid fields error",
	}

	if !err.isInvalidFields() {
		t.Errorf("Expected isInvalidFields to return true")
	}

	// Test with other type
	err2 := &PxError{
		Type:    "NOT_FOUND",
		Message: "Not found error",
	}

	if err2.isInvalidFields() {
		t.Errorf("Expected isInvalidFields to return false for NOT_FOUND type")
	}
}

func TestPxError_isNotFound(t *testing.T) {
	// Test with NOT_FOUND type
	err := &PxError{
		Type:    "NOT_FOUND",
		Message: "Resource not found",
	}

	if !err.isNotFound() {
		t.Errorf("Expected isNotFound to return true")
	}

	// Test with other type
	err2 := &PxError{
		Type:    "INVALID_FIELDS",
		Message: "Invalid fields",
	}

	if err2.isNotFound() {
		t.Errorf("Expected isNotFound to return false for INVALID_FIELDS type")
	}
}

func TestPxError_Error(t *testing.T) {
	// Test error without fields
	err := &PxError{
		Type:    "NOT_FOUND",
		Message: "Resource not found",
	}

	if err.Error() != "Resource not found" {
		t.Errorf("Expected 'Resource not found', got '%s'", err.Error())
	}

	// Test error with fields
	errWithFields := &PxError{
		Type:    "INVALID_FIELDS",
		Message: "Validation failed",
		Fields: []PxInvalidField{
			{Name: "Name", Reason: "required", Message: "Name is required"},
			{Name: "Email", Reason: "format", Message: "Invalid email format"},
		},
	}

	errorMsg := errWithFields.Error()
	if !strings.Contains(errorMsg, "Validation failed") {
		t.Errorf("Expected error message to contain 'Validation failed', got '%s'", errorMsg)
	}
	if !strings.Contains(errorMsg, "Name") || !strings.Contains(errorMsg, "Email") {
		t.Errorf("Expected error message to contain field names, got '%s'", errorMsg)
	}
}

func TestNewPxError(t *testing.T) {
	// Test with valid JSON error response
	jsonError := `{"Type":"NOT_FOUND","Message":"Resource not found","Status":404}`
	reader := strings.NewReader(jsonError)

	err := NewPxError(reader, 404, "ADR/Adresse/123")

	if err.Status != 404 {
		t.Errorf("Expected status 404, got %d", err.Status)
	}

	if err.Type != "NOT_FOUND" {
		t.Errorf("Expected type 'NOT_FOUND', got '%s'", err.Type)
	}

	if err.Endpoint != "ADR/Adresse/123" {
		t.Errorf("Expected endpoint 'ADR/Adresse/123', got '%s'", err.Endpoint)
	}

	// Test with nil reader
	nilErr := NewPxError(nil, 500, "PRO/Login")

	if nilErr.Status != 500 {
		t.Errorf("Expected status 500, got %d", nilErr.Status)
	}

	if nilErr.Message == "" {
		t.Errorf("Expected default message for nil reader")
	}

	// Test with empty message - should get default
	emptyReader := strings.NewReader(`{"Type":"ERROR"}`)
	emptyErr := NewPxError(emptyReader, 400, "TEST/Endpoint")

	if emptyErr.Message == "" {
		t.Errorf("Expected default message when API doesn't provide one")
	}

	if !strings.Contains(emptyErr.Message, "400") {
		t.Errorf("Expected default message to contain status code, got '%s'", emptyErr.Message)
	}
}

func TestNewPxError_WithInvalidFields(t *testing.T) {
	// Test with INVALID_FIELDS error
	jsonError := `{
		"Type":"INVALID_FIELDS",
		"Message":"Validation failed",
		"Fields":[
			{"Name":"AdressNr","Reason":"required","Message":"AdressNr is required"},
			{"Name":"Name","Reason":"length","Message":"Name must be at least 3 characters"}
		]
	}`
	reader := strings.NewReader(jsonError)

	err := NewPxError(reader, 400, "ADR/Adresse")

	if err.Type != "INVALID_FIELDS" {
		t.Errorf("Expected type 'INVALID_FIELDS', got '%s'", err.Type)
	}

	if len(err.Fields) != 2 {
		t.Errorf("Expected 2 invalid fields, got %d", len(err.Fields))
	}

	if err.Fields[0].Name != "AdressNr" {
		t.Errorf("Expected first field name 'AdressNr', got '%s'", err.Fields[0].Name)
	}

	if err.Fields[1].Reason != "length" {
		t.Errorf("Expected second field reason 'length', got '%s'", err.Fields[1].Reason)
	}
}

func TestNewPxError_NoStatusNoEndpoint(t *testing.T) {
	// Test with no status and no endpoint
	err := NewPxError(nil, 0, "")

	if err.Message != "unknown error" {
		t.Errorf("Expected 'unknown error' message, got '%s'", err.Message)
	}
}
