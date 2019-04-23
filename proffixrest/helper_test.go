package proffixrest

import (
	"context"
	"github.com/spf13/cast"
	"strings"
	"testing"
)

func Test_errorFormatterPx(t *testing.T) {

	//New Context
	ctx := context.Background()

	pxrest, _ := ConnectTest(ctx, []string{"ADR"})

	_, _, status, err := pxrest.Get(ctx, "ADR/Adresse/12345678", nil)

	pxrest.Logout(ctx)

	if !strings.Contains(err.Error(), cast.ToString(status)) {
		t.Errorf("Error should contain Statuscode %v", status)
	}

	var notfound = "NOT_FOUND"

	if !strings.Contains(err.Error(), notfound) {
		t.Errorf("Error should contain Type %v", notfound)
	}

}
