package proffixrest

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strconv"
)

// File Write Request for PROFFIX REST-API
// Accepts Context, Endpoint and []Byte as Input
// Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) File(ctx context.Context, filename string, data []byte) (io.ReadCloser, http.Header, int, error) {
	err := c.Login(ctx)

	if err != nil {
		return nil, nil, 0, err
	}

	// Define endpoint
	var endpoint = "PRO/Datei"

	// Set filename

	params := url.Values{}
	if filename != "" {
		params.Set("filename", filename)
	}
	request, header, statuscode, err := c.request(ctx, "POST", endpoint, params, true, data)

	//If Log enabled in options log data
	logDebug(ctx, c, fmt.Sprintf("Sent data in POST-Request: %v", data))

	if err != nil {
		return request, header, statuscode, err
	} else {
		return request, header, statuscode, nil
	}
}

// File Read Request for PROFFIX REST-API
// Accepts Context, DateiNr and Params as Input
// Returns File as io.ReadCloser,Filename string,ContentType string,ContentLength int,error
func (c *Client) GetFile(ctx context.Context, dateinr string, params url.Values) (rc io.ReadCloser, FileName string, ContentType string, ContentLength int, err error) {

	// Build query for getting download URL of List
	resp, headers, status, err := c.Get(ctx, "PRO/Datei/"+dateinr, params)

	// If err not nil or status not 200
	if err != nil || status != 200 {
		logDebug(ctx, c, fmt.Sprintf("Error on create list: %v, PxSession-ID: %v", err, c.GetPxSessionId()))
		return resp, "", "", 0, err
	}

	ContentType = headers.Get("Content-Type")
	ContentLength, _ = strconv.Atoi(headers.Get("Content-Length"))
	ContentDisposition := headers.Get("Content-Disposition")

	_, paramsCD, errCD := mime.ParseMediaType(ContentDisposition)
	if errCD == nil {
		FileName = paramsCD["filename"]
	}

	// If Log enabled log URL
	logDebug(ctx, c, fmt.Sprintf("Downloaded File with Content-Length: %v, PxSession-ID: %v", ContentLength, c.GetPxSessionId()))

	return resp, FileName, ContentType, ContentLength, err

}
