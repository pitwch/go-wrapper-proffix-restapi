package proffixrest

import (
	"context"
	"os"
	"testing"
)

func TestClient_CheckAPI(t *testing.T) {

	ctx := context.Background()

	// Connect
	pxrest, err := ConnectTest([]string{})

	// Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Connect. Got '%v'", err)
	}

	// Check CheckAPI
	err = pxrest.CheckAPI(ctx, os.Getenv("PXDEMO_KEY"))

	// Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for CheckAPI. Got '%v'", err)
	}

}
