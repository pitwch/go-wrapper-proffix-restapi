package proffixrest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) GetList(ctx context.Context, listname string, body interface{}) (io.ReadCloser, http.Header, int, error) {

	//Build query for getting download URL of List
	resp, headers, status, err := c.Post(ctx, "PRO/Liste/"+listname+"/generieren", body)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp)
	res := buf.String()

	//If err not nil or status not 201
	if err != nil || status != 201 {
		return resp, headers, status, err
	}

	downloadLocation := headers.Get("Location")

	//If Log enabled log URL
	logDebug(ctx, c, fmt.Sprintf("Got Download URL from PROFFIX REST-API: %v, PxSession-ID: %v", downloadLocation, PxSessionId))

	//Remove Path from Response -> replace with empty string
	cleanLocation := strings.Replace(downloadLocation, c.restURL.String(), "", -1)

	downloadFile, headersDownload, statusDownload, err := c.Get(ctx, cleanLocation, nil)

	return downloadFile, headersDownload, statusDownload, err

}
