package proffixrest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// SyncBatchData is a placeholder for batch requests
type SyncBatchData map[string]interface{}

// SyncBatch automatically POST/PUT items based on Keyfield in Body
// Params:
//
//	ctx				context.Context	context
//	endpoint		string			the endpoint on which batch should be used
//	keyfield		string			the keyfield from endpoint
//	removeKeyfield	bool			if true removes the keyfield on POST Request
//	data			[]byte			JSON Request as []byte
//
// Returns:
//
//	created			[]string		Keynrs of created entries (POST)
//	updated			[]string		Keynrs of updated entries (PUT)
//	failed			[]string		Keynrs of failed entries
//	errors			[]string		Errors as string from atch request
//	total			int				Total requests
//	err				error			General errors
func (c *Client) SyncBatch(ctx context.Context, endpoint string, keyfield string, removeKeyfield bool, data []byte) (created []string, updated []string, failed []string, errors []string, total int, err error) {

	var datas []SyncBatchData

	err = json.Unmarshal(data, &datas)
	if err != nil {
		return nil, nil, nil, nil, 0, err
	}

	for i := range datas {

		// Get Key from Map
		key := fmt.Sprintf("%v", datas[i][keyfield])

		var (
			statusGet int
			getResp   io.Reader
		)

		// If Keyfield is empty / missing (saves a GET Request...)
		if datas[i][keyfield] == "" {
			statusGet = 404
		} else {
			getResp, _, statusGet, _ = c.Get(ctx, endpoint+"/"+key, nil)
		}

		switch statusGet {
		case 404:
			// If Item not found -> create / post it with extracted keyfield
			if removeKeyfield {
				delete(datas[i], keyfield)
			}

			_, headers, status, err := c.Post(ctx, endpoint, datas[i])
			if status == 201 {
				// Append to created
				created = append(created, ConvertLocationToID(headers))
			} else {
				// Append to failed
				failed = append(failed, key)
				errors = append(errors, fmt.Sprintf("%v", err))
			}

		case 200:
			// If Item found -> update / put new values
			res := ""
			resp, _, status, err := c.Put(ctx, endpoint+"/"+key, datas[i])
			if resp != nil {
				// Buffer decode for plain text response
				buf := new(bytes.Buffer)
				_, _ = buf.ReadFrom(resp)
				res = buf.String()
			}

			if status == 204 {
				// Append to updated
				updated = append(updated, key)
			} else {
				// Append to failed
				failed = append(failed, key)
				errors = append(errors, fmt.Sprintf("%v %v", err, res))
			}

		default:
			// Unexpected status; append to failed
			res := ""
			if getResp != nil {
				// Buffer decode for plain text response
				buf := new(bytes.Buffer)
				_, _ = buf.ReadFrom(getResp)
				res = buf.String()
			}
			failed = append(failed, key)
			errors = append(errors, fmt.Sprintf("%v %v", err, res))
		}

	}

	return created, updated, failed, errors, len(datas), err
}
