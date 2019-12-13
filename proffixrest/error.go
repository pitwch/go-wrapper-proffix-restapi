package proffixrest

import (
	"context"
	"fmt"
	"io"
)

// proffixErrorStruct defines the Error Struct of PROFFIX REST-API
type proffixErrorStruct struct {
	Type    string `json:"Type"`
	Message string `json:"Message"`
	Fields  []struct {
		Reason  string `json:"Reason"`
		Name    string `json:"Name"`
		Message string `json:"Message"`
	} `json:"Fields"`
}

type ErrorStruct struct {
	Message string   `json:"Message"`
	Fields  []string `json:"Fields"`
}

// errorFormatterPx does format the errors and handles emergency logout
func errorFormatterPx(ctx context.Context, c *Client, statuscode int, request io.Reader) (err error) {
	errstr, _ := ReaderToString(request)
	// No Autologout if deactivated in options
	if c.option.Autologout {
		// Error 404 is soft so no logout
		if statuscode != 404 {
			// Do Forced Logout...
			logoutstatus, _ := c.Logout(ctx)
			logDebug(ctx, c, fmt.Sprintf("Logout after Error with Status: %v", logoutstatus))

		}
	}

	return fmt.Errorf("%v", errstr)
}

// Todo: ErrorBeautifier beautifies Errors from PROFFIX
/*func ErrorBeautifier(err error) (){

}*/
