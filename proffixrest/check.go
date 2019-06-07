package proffixrest

import (
	"context"
	"fmt"
	"net/url"
)

// CheckAPI checks the PROFFIX REST API
// Accepts client
// Returns err
func (c *Client) CheckApi(ctx context.Context) (err error) {

	//Set timeout to 10
	c.option.Timeout = 10

	//Set Module ADR
	c.Module = []string{}

	//Check Login
	err = c.Login(ctx)
	if err != nil {
		return err
	}

	//Check if PxSessionId is availabler
	if c.GetPxSessionId() == "" {
		return fmt.Errorf("Keine PxSessionId gefunden")
	}

	//Check if GET ADR/Adresse works

	//Set params for minimal request
	params := url.Values{}
	params.Add("Limit", "1")
	params.Add("Fields", "AdressNr")
	_, _, statuscode, err := c.Get(ctx, "ADR/Adresse", params)

	//To make sure: Logout
	_, _ = c.Logout(ctx)

	//Check if error
	if err != nil {
		return err
	}

	//Also check statuscode
	if statuscode != 200 {
		return fmt.Errorf("Falscher Statuscode bei Testabfrage auf ADR/Adresse: %v", statuscode)
	}

	return err

}
