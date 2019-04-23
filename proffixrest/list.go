package proffixrest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (c *Client) GetList(ctx context.Context, listenr int, body interface{}) (io.ReadCloser, http.Header, int, error) {

	//Build query for getting download URL of List
	resp, headers, status, err := c.Post(ctx, "PRO/Liste/"+strconv.Itoa(listenr)+"/generieren", body)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp)
	res := buf.String()

	//If err not nil or status not 201
	if err != nil || status != 201 {
		logDebug(ctx, c, fmt.Sprintf("Error on create list: %v, PxSession-ID: %v", res, PxSessionId))
		return resp, headers, status, err
	}

	downloadLocation := headers.Get("Location")

	//If Log enabled log URL
	logDebug(ctx, c, fmt.Sprintf("Got Download URL from PROFFIX REST-API: %v, PxSession-ID: %v", downloadLocation, PxSessionId))

	//Remove Path from Response -> replace with empty string
	cleanLocation := strings.Replace(downloadLocation, c.restURL.String(), "", -1)

	downloadFile, headersDownload, statusDownload, err := c.Get(ctx, cleanLocation, nil)

	//If Log enabled log URL
	logDebug(ctx, c, fmt.Sprintf("Downloaded File from '%v' with Content-Length: %v, PxSession-ID: %v", cleanLocation, headersDownload.Get("Content-Length"), PxSessionId))

	return downloadFile, headersDownload, statusDownload, err

}
