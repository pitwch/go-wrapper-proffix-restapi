package proffixrest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// File Write Request for PROFFIX REST-API
// Accepts Context, Endpoint and []Byte as Input
// Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) File(ctx context.Context, filename string, data []byte) (io.ReadCloser, http.Header, int, error) {
	err := c.Login(ctx)

	if !err.(*PxError).isNull() {
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

	if !err.(*PxError).isNull() {
		return request, header, statuscode, err
	} else {
		return request, header, statuscode, nil
	}
}

// File Read Request for PROFFIX REST-API
// Accepts Context and DateiNr as Input
// Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) GetFile(ctx context.Context, dateinr string) (io.ReadCloser, http.Header, int, error) {

	// Build query for getting download URL of List
	resp, headers, status, err := c.Get(ctx, "PRO/Datei/"+dateinr, nil)

	// If err not nil or status not 201
	if err != nil || status != 201 {
		logDebug(ctx, c, fmt.Sprintf("Error on create list: %v, PxSession-ID: %v", err.(*PxError).Message, PxSessionId))
		return resp, headers, status, err
	}

	downloadLocation := headers.Get("Location")

	// If Log enabled log URL
	logDebug(ctx, c, fmt.Sprintf("Got Download URL from PROFFIX REST-API: %v, PxSession-ID: %v", downloadLocation, PxSessionId))

	// Remove Path from Response -> replace with empty string
	cleanLocation := strings.Replace(downloadLocation, c.restURL.String(), "", -1)

	downloadFile, headersDownload, statusDownload, err := c.Get(ctx, cleanLocation, nil)

	if headersDownload == nil {
		headersDownload = http.Header{}
	}

	// If Log enabled log URL
	logDebug(ctx, c, fmt.Sprintf("Downloaded File from '%v' with Content-Length: %v, PxSession-ID: %v", cleanLocation, headersDownload.Get("Content-Length"), PxSessionId))

	return downloadFile, headersDownload, statusDownload, err

}
