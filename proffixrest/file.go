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

// File uploads binary data to the PROFFIX REST-API.
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

	// If Log enabled in options log data
	logDebug(ctx, c, fmt.Sprintf("Sent data in POST-Request: %v", data))

	if err != nil {
		return request, header, statuscode, err
	}
	return request, header, statuscode, nil
}

// GetFile retrieves a file from the PROFFIX REST-API by its identifier.
// Accepts Context, DateiNr and Params as Input
// Returns File as io.ReadCloser,Filename string,ContentType string,ContentLength int,error
func (c *Client) GetFile(ctx context.Context, dateinr string, params url.Values) (rc io.ReadCloser, fileName string, contentType string, contentLength int, err error) {

	// Build query for getting download URL of List
	resp, headers, status, err := c.Get(ctx, "PRO/Datei/"+dateinr, params)

	if err != nil || status != 200 {
		logDebug(ctx, c, fmt.Sprintf("Error on fetch file: %v, PxSession-ID: %v", err, c.GetPxSessionID()))
		return resp, "", "", 0, err
	}

	contentType = headers.Get("Content-Type")
	contentLength, _ = strconv.Atoi(headers.Get("Content-Length")) // Ignore error, default to 0
	ContentDisposition := headers.Get("Content-Disposition")

	_, paramsCD, errCD := mime.ParseMediaType(ContentDisposition)
	if errCD == nil {
		fileName = paramsCD["filename"]
	}

	// If Log enabled log URL
	logDebug(ctx, c, fmt.Sprintf("Downloaded File with Content-Length: %v, PxSession-ID: %v", contentLength, c.GetPxSessionID()))

	return resp, fileName, contentType, contentLength, err

}
