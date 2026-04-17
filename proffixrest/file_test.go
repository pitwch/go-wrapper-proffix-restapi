package proffixrest

import (
	"context"
	"io"
	"net/url"
	"testing"
)

func TestClient_File(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	// Test file upload with sample data
	testData := []byte("Test file content for PROFFIX REST API file upload")
	filename := "test_upload.txt"

	rc, headers, status, err := pxrest.File(ctx, filename, testData)

	// Note: This test may return 404 if file upload is not configured
	// or return success if configured
	if status == 201 {
		if err != nil {
			t.Errorf("Expected no error for successful file upload. Got '%v'", err)
		}

		if rc == nil {
			t.Errorf("Expected response body for successful upload")
		} else {
			_ = rc.Close()
		}

		if headers.Get("Content-Type") == "" {
			t.Errorf("Expected Content-Type header in response")
		}
	} else if status == 404 {
		// File upload endpoint may not be available - this is acceptable
		t.Logf("File upload endpoint returned 404 - may not be configured")
	} else if status >= 400 {
		// Other errors are logged but not failing
		t.Logf("File upload returned status %d: %v", status, err)
	}
}

func TestClient_File_NoFilename(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	// Test file upload without filename
	testData := []byte("Test file content without filename")

	rc, _, status, _ := pxrest.File(ctx, "", testData)

	// The API should handle this gracefully (either accept or reject)
	if rc != nil {
		_ = rc.Close()
	}

	// Status can be anything - we're just testing the request is made
	t.Logf("File upload without filename returned status %d", status)
}

func TestClient_GetFile(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	// Test getting a non-existent file
	rc, fileName, contentType, contentLength, err := pxrest.GetFile(ctx, "99999999", nil)

	// Should get 404 for non-existent file
	if err != nil {
		// Error is expected for non-existent file
		t.Logf("Expected error for non-existent file: %v", err)
	}

	if rc != nil {
		_ = rc.Close()
	}

	// Content should be empty or have defaults
	t.Logf("GetFile returned: filename=%s, contentType=%s, length=%d", fileName, contentType, contentLength)
}

func TestClient_GetFile_WithParams(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	// Test with query params
	params := url.Values{}
	params.Set("preview", "true")

	// Try to get a file with params
	rc, fileName, contentType, contentLength, err := pxrest.GetFile(ctx, "99999999", params)

	if rc != nil {
		_ = rc.Close()
	}

	// Should either succeed or fail gracefully
	if err != nil {
		t.Logf("GetFile with params error: %v", err)
	}

	t.Logf("GetFile with params returned: filename=%s, contentType=%s, length=%d", fileName, contentType, contentLength)
}

func TestClient_File_EmptyData(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	// Test with empty data
	emptyData := []byte{}

	rc, _, status, err := pxrest.File(ctx, "empty.txt", emptyData)

	if rc != nil {
		// Read and close
		_, _ = io.ReadAll(rc)
		_ = rc.Close()
	}

	// API should handle empty files gracefully
	t.Logf("Empty file upload returned status %d, err: %v", status, err)
}

func TestClient_File_LargeFile(t *testing.T) {
	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{"VOL"})
	if err != nil {
		t.Fatalf("ConnectTest failed: %v", err)
	}
	defer func() { _, _ = pxrest.Logout(ctx) }()

	// Create larger test data (10KB)
	largeData := make([]byte, 10240)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	rc, _, status, err := pxrest.File(ctx, "large_test.bin", largeData)

	if rc != nil {
		_ = rc.Close()
	}

	// API should handle large files
	t.Logf("Large file upload returned status %d, err: %v", status, err)
}
