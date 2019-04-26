package proffixrest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type SyncBatchData map[string]interface{}

// SyncBatch automatically POST/PUT items based on Keyfield in Body
//		ctx	context.Context	context
//		endpoint	string	the endpoint on which batch should be used
//		keyfield	string	the keyfield from endpoint
//		removeKeyfield	bool	if true removes the keyfield on POST Request
//		data	[]byte	JSON Request as []byte
// Returns result (keyfield,status,err) , total and error
func (c *Client) SyncBatch(ctx context.Context, endpoint string, keyfield string, removeKeyfield bool, data []byte) (created []string, updated []string, failed []string, errors []string, total int, err error) {

	var datas []SyncBatchData

	err = json.Unmarshal([]byte(data), &datas)
	if err != nil {
		log.Printf("Error on Unmarshal: %v", err)
	}

	for i := range datas {

		//Get Key from Map
		key := fmt.Sprintf("%v", datas[i][keyfield])

		statusGet := 0
		var getResp io.Reader
		//If Keyfield is empty / missing (saves a GET Request...)
		if datas[i][keyfield] == "" {
			statusGet = 404
		} else {
			getResp, _, statusGet, _ = c.Get(ctx, endpoint+"/"+key, nil)

		}

		//If Item not found -> create / post it with extracted keyfield
		if statusGet == 404 {
			res := ""

			//If flag removeKeyfield is set -> remove Keyfield from Post - Request
			if removeKeyfield {
				delete(datas[i], keyfield)
			}

			resp, headers, status, err := c.Post(ctx, endpoint, datas[i])

			x, err := json.Marshal(datas[i])
			log.Println(string(x))

			if resp != nil {
				//Buffer decode for plain text response
				buf := new(bytes.Buffer)
				buf.ReadFrom(resp)
				res = buf.String()
			}

			if status == 201 {

				//Append to created
				created = append(created, ConvertLocationToID(headers))

			} else {
				//Append to failed
				failed = append(failed, key)

				errors = append(errors, fmt.Sprintf("%v %v", err, res))
			}

			//If Item found -> update / put new values
		} else if statusGet == 200 {
			res := ""
			resp, _, status, err := c.Put(ctx, endpoint+"/"+key, datas[i])
			if resp != nil {
				//Buffer decode for plain text response
				buf := new(bytes.Buffer)
				buf.ReadFrom(resp)
				res = buf.String()
			}

			if status == 204 {

				//Append to updated
				updated = append(updated, key)

			} else {
				//Append to failed
				failed = append(failed, key)

				errors = append(errors, fmt.Sprintf("%v %v", err, res))
			}
		} else {
			res := ""
			if getResp != nil {
				//Buffer decode for plain text response
				buf := new(bytes.Buffer)
				buf.ReadFrom(getResp)
				res = buf.String()
			}

			//Append to failed
			failed = append(failed, key)

			errors = append(errors, fmt.Sprintf("%v %v", err, res))
		}

	}

	return created, updated, failed, errors, len(datas), err
}
