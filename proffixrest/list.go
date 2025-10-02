package proffixrest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// GetList generates a list on the PROFFIX REST-API and downloads its file representation.
func (c *Client) GetList(ctx context.Context, listenr int, body interface{}) (io.ReadCloser, http.Header, int, error) {

	// Build query for getting download URL of List
	resp, headers, status, err := c.Post(ctx, "PRO/Liste/"+strconv.Itoa(listenr)+"/generieren", body)

	// If err not nil or status not 201
	if err != nil || status != 201 {
		logDebug(ctx, c, fmt.Sprintf("Error on create list: %v, PxSession-ID: %v", err, c.GetPxSessionID()))
		return resp, headers, status, err
	}

	// Build Download Uri - was change in Px Rest-API 4.43 / 4.42. This fix works for all...
	dateiNr := ConvertLocationToID(headers)
	downloadURI := "PRO/Datei/" + dateiNr

	// If Log enabled log URL
	logDebug(ctx, c, fmt.Sprintf("Got Download URL from PROFFIX REST-API: %v, PxSession-ID: %v", downloadURI, c.GetPxSessionID()))

	downloadFile, headersDownload, statusDownload, err := c.Get(ctx, downloadURI, nil)

	if headersDownload == nil {
		headersDownload = http.Header{}
	}

	// If Log enabled log URL
	logDebug(ctx, c, fmt.Sprintf("Downloaded File from '%v' with Content-Length: %v, PxSession-ID: %v", downloadURI, headersDownload.Get("Content-Length"), c.GetPxSessionID()))

	return downloadFile, headersDownload, statusDownload, err

}
