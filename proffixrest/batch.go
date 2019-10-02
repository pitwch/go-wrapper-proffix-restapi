package proffixrest

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strconv"
)

// 	GetBatch automatically paginates all possible queries
//	Params:
//		ctx				context.Context	context
//		endpoint		string			the endpoint on which batch should be used
//		params			url.Values		the params for this request
//		batchsize		int				the size of a batch (higher is usually faster)
//
//	Returns:
//		result			[]byte			results as []byte
//		total			int				Total found entries
//		err				error			General errors
func (c *Client) GetBatch(ctx context.Context, endpoint string, params url.Values, batchsize int) (result []byte, total int, err error) {

	// Create collector for collecting requests
	var collector []byte

	// If batchsize not set in function / is 0
	if batchsize == 0 {
		//Set it to default
		batchsize = c.option.Batchsize
	}

	// If params were set to nil set to empty url.Values{}
	if params == nil {
		params = url.Values{}
	}

	// Create Setting Query for setting up the sync and getting the FilteredCount Header

	// Store original params in paramquery
	paramquery := params

	// Delete Limit from params in case user defined it
	paramquery.Del("Limit")

	// Add Limit / Batchsize to params
	paramquery.Add("Limit", strconv.Itoa(batchsize))
	// Query Endpoint for results
	rc, header, status, err := c.Get(ctx, endpoint, params)

	if err != nil {
		return nil, 0, err
	}

	// Read from rc into settingResp
	settingResp, err := ioutil.ReadAll(rc)

	// Put setting query into collector
	collector = append(collector, settingResp[:]...)

	// If Status not 200 log Error message
	if status != 200 {
		return nil, 0, err
	}
	// Get total available objects
	totalEntries := GetFiltererCount(header)

	// Unmarshall to interface so we can count elements in setting query
	var jsonObjs interface{}
	err = json.Unmarshal([]byte(collector), &jsonObjs)

	firstQueryCount := len(jsonObjs.([]interface{}))

	// Add firstQueryCount to the totalEntriesCount so we can save queries
	totalEntriesCount := firstQueryCount

	// If total available objects are already all requested --> return and exit
	if totalEntriesCount < totalEntries {

		// Loop over requests until we have all available objects
		for totalEntriesCount < totalEntries {

			// Reset the params to original
			paramquery := url.Values{}
			paramquery = params

			// To be sure; remove old limit + offset
			paramquery.Del("Limit")
			paramquery.Del("Offset")

			// Add new limit + offset
			paramquery.Add("Limit", strconv.Itoa(batchsize))
			paramquery.Add("Offset", strconv.Itoa(totalEntriesCount))

			// Fire request with new limit + offset
			bc, _, _, _ := c.Get(ctx, endpoint, paramquery)

			// Read from bc into temporary batchResp
			batchResp, err := ioutil.ReadAll(bc)

			collector = append(collector, batchResp[:]...)
			if err != nil {
				panic(err)
			}

			// If Status not 200 log Error message
			if status != 200 {
				panic(string(collector))
			}

			// Unmarshall to interface so we can count elements in totalEntriesCount query
			var jsonObjs2 interface{}
			_ = json.Unmarshal([]byte(batchResp), &jsonObjs2)
			totalEntriesCount += len(jsonObjs2.([]interface{}))
		}

		// Clean the collector, replace JSON Tag to get one big JSON array
		cleanCollector := bytes.Replace(collector, []byte("]["), []byte(","), -1)

		return cleanCollector, totalEntriesCount, err

		// If count of elements in first query is bigger or same than FilteredCount return
	} else {
		// Return
		return collector, totalEntriesCount, err
	}
}
