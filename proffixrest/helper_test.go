package proffixrest

import (
	"context"
	"testing"
)

func Test_errorFormatterPx(t *testing.T) {

	//New Context
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
