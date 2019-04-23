package proffixrest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

//Sync Request for PROFFIX REST-API
//Accepts Context, Check, endpointCheck, endpointWrite and Data as Input
//Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Sync(ctx context.Context, endpointCheck string, endpointWrite string, keyfield string, data interface{}) (io.ReadCloser, http.Header, int, error) {
	err := c.Login(ctx)
	if err != nil {
		return nil, nil, 0, err
	}

	//GET Request for checking if item already exist
	_, _, checkstatus, _ := c.request(ctx, "GET", endpointCheck+"/"+keyfield, nil, nil)

	//Define Vars
	var request io.ReadCloser
	var header http.Header
	var statuscode int

	//If Status 200 -> item already exist -> PUT
	if checkstatus == 200 {
		request, header, statuscode, err = c.request(ctx, "PUT", endpointWrite+"/"+keyfield, url.Values{}, data)
		//If Status not 200 -> item doesnt exist -> POST
	} else {
		request, header, statuscode, err = c.request(ctx, "POST", endpointWrite, url.Values{}, data)
	}

	//If Log enabled in options log data
	logDebug(ctx, c, fmt.Sprintf("Sent data in POST-Request: %v", data))

	//If Login is invalid - try again
	if statuscode == 401 {
		//Set PxSessionId to empty string
		PxSessionId = ""

		//Get new pxsessionid and write to var
		err = c.Login(ctx)

		//Repeat Request with new SessionId
		//If Status 200 -> item already exist -> PUT
		if checkstatus == 200 {
			//a := tokenGenerator()
			request, header, statuscode, err = c.request(ctx, "PUT", endpointWrite, url.Values{}, data)
			//If Status not 200 -> item doesnt exist -> POST
		} else {
			request, header, statuscode, err = c.request(ctx, "POST", endpointWrite, url.Values{}, data)
		}

		return request, header, statuscode, err
	}

	//If Statuscode not 201
	if statuscode != 201 && statuscode != 204 {
		return request, header, statuscode, errorFormatterPx(ctx, c, statuscode, request)
	}

	return request, header, statuscode, err
}
