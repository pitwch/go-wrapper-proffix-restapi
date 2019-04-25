package proffixrest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type SyncBatchData map[string]interface{}

// SyncBatch automatically POST/PUT items based on Keyfield in Body
// Accepts context, endpoint, keyfield and data as JSON Array
// Returns result (keyfield,status,err) , total and error
func (c *Client) SyncBatch(ctx context.Context, endpoint string, keyfield string, data []byte) (result [][]string, total int, err error) {

	var datas []SyncBatchData

	err = json.Unmarshal([]byte(data), &datas)
	if err != nil {
		log.Printf("Error on Unmarshal: %v", err)
	}

	for i := range datas {

		//Get Key from Map
		key := fmt.Sprintf("%v", datas[i][keyfield])

		statusGet := 0

		//If Keyfield is empty / missing (saves a GET Request...)
		if datas[i][keyfield] == nil {
			statusGet = 404
		} else {
			log.Printf("TEST 404 %v", datas[i][keyfield])
			log.Println(key)
			_, _, statusGet, _ = c.Get(ctx, endpoint+"/"+key, nil)
		}
		//If Item not found -> create / post it with extracted keyfield
		if statusGet == 404 {
			res := ""
			delete(datas[i], keyfield)
			resp, headers, status, err := c.Post(ctx, endpoint, datas[i])

			if resp != nil {
				//Buffer decode for plain text response
				buf := new(bytes.Buffer)
				buf.ReadFrom(resp)
				res = buf.String()
			}
			//Append to result
			result = append(result, []string{ConvertLocationToID(headers), strconv.Itoa(status), res, fmt.Sprintf("%v", err)})
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
			result = append(result, []string{key, strconv.Itoa(status), res, fmt.Sprintf("%v", err)})

		}
	}

	return result, len(datas), err
}
