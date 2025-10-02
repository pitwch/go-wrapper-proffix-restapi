package proffixrest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (c *Client) GetList(ctx context.Context, listenr int, body interface{}) (io.ReadCloser, http.Header, int, error) {

	//Build query for getting download URL of List
	resp, headers, status, err := c.Post(ctx, "PRO/Liste/"+strconv.Itoa(listenr)+"/generieren", body)

	//If err not nil or status not 201
	if err != nil || status != 201 {
		logDebug(ctx, c, fmt.Sprintf("Error on create list: %v, PxSession-ID: %v", err, c.GetPxSessionId()))
		return resp, headers, status, err
	}

	// Build Download Uri - was change in Px Rest-API 4.43 / 4.42. This fix works for all...
	dateiNr := ConvertLocationToID(headers)
	downloadUri := "PRO/Datei/" + dateiNr

	//If Log enabled log URL
	logDebug(ctx, c, fmt.Sprintf("Got Download URL from PROFFIX REST-API: %v, PxSession-ID: %v", downloadUri, c.GetPxSessionId()))

	downloadFile, headersDownload, statusDownload, err := c.Get(ctx, downloadUri, nil)

	if headersDownload == nil {
		headersDownload = http.Header{}
	}

	//If Log enabled log URL
	logDebug(ctx, c, fmt.Sprintf("Downloaded File from '%v' with Content-Length: %v, PxSession-ID: %v", downloadUri, headersDownload.Get("Content-Length"), c.GetPxSessionId()))

	return downloadFile, headersDownload, statusDownload, err

}
