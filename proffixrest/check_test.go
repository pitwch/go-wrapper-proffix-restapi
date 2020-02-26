package proffixrest

import (
	"context"
	"os"
	"testing"
)

func TestClient_CheckApi(t *testing.T) {

	ctx := context.Background()

	//Connect
	pxrest, err := ConnectTest([]string{"ADR"})

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Connect. Got '%v'", err)
	}

	//Check CheckAPI
	err = pxrest.CheckApi(ctx, os.Getenv("PXDEMO_KEY"))

	//Check error. Should be nil
	if !err.(*PxError).isNull() {
		t.Errorf("Expected no error for CheckApi. Got '%v'", err)
	}

}
