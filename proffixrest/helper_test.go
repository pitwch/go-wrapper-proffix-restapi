package proffixrest

import (
	"context"
	"strings"
	"testing"
)

func Test_errorFormatterPx(t *testing.T) {

	//New Context
	ctx := context.Background()

	pxrest, _ := ConnectTest(ctx, []string{"ADR"})

	_, _, status, err := pxrest.Get(ctx, "ADR/Adresse/12345678", nil)

	_, _ = pxrest.Logout(ctx)

	if status != 404 {
		t.Errorf("Statuscode should be 404: Statuscode %v", status)
	}

	var notfound = "NOT_FOUND"

	if !strings.Contains(err.Error(), notfound) {
		t.Errorf("Error should contain Type %v", notfound)
	}

}
