package proffixrest

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/spf13/cast"
	"io/ioutil"
	"net/url"
)

// GetBatch automatically paginates until all possible queries were received.
// Accepts context, endpoint, params and batchsize
// Returns []byte,total, err
func (c *Client) GetBatch(ctx context.Context, endpoint string, params url.Values, batchsize int) (result []byte, total int, err error) {

	//Create collector for collecting requests
	var collector []byte

	//If batchsize not set in function / is 0
	if batchsize == 0 {
		//Set it to default
		batchsize = c.option.Batchsize
	}

	//If params were set to nil set to empty url.Values{}
	if params == nil {
		params = url.Values{}
	}

	//Create Setting Query for setting up the sync and getting the FilteredCount Header

	//Store original params in paramquery
	paramquery := params

	//Delete Limit from params in case user defined it
	paramquery.Del("Limit")

	//Add Limit / Batchsize to params
	paramquery.Add("Limit", cast.ToString(batchsize))

	//Query Endpoint for results
	rc, header, status, err := c.Get(ctx, endpoint, params)

	if err != nil {
		return nil, 0, err
	}

	//Read from rc into settingResp
	settingResp, err := ioutil.ReadAll(rc)

	//Put setting query into collector
	collector = append(collector, settingResp[:]...)

	//If Status not 200 log Error message
	if status != 200 {
		return nil, 0, err
	}
	//Get total available objects
	totalEntries := GetFiltererCount(header)

	//Unmarshall to interface so we can count elements in setting query
	var jsonObjs interface{}
	err = json.Unmarshal([]byte(collector), &jsonObjs)

	firstQueryCount := len(jsonObjs.([]interface{}))

	//Add firstQueryCount to the totalEntriesCount so we can save queries
	totalEntriesCount := firstQueryCount

	//If total available objects are already all requested --> return and exit
	if totalEntriesCount < totalEntries {

		//Loop over requests until we have all available objects
		for totalEntriesCount < totalEntries {

			//Reset the params to original
			paramquery := url.Values{}
			paramquery = params

			//To be sure; remove old limit + offset
			paramquery.Del("Limit")
			paramquery.Del("Offset")

			//Add new limit + offset
			paramquery.Add("Limit", cast.ToString(batchsize))
			paramquery.Add("Offset", cast.ToString(totalEntriesCount))

			//Fire request with new limit + offset
			bc, _, _, _ := c.Get(ctx, endpoint, paramquery)

			//Read from bc into temporary batchResp
			batchResp, err := ioutil.ReadAll(bc)

			collector = append(collector, batchResp[:]...)
			if err != nil {
				panic(err)
			}

			//If Status not 200 log Error message
			if status != 200 {
				panic(string(collector))
			}

			//Unmarshall to interface so we can count elements in totalEntriesCount query
			var jsonObjs2 interface{}
			json.Unmarshal([]byte(batchResp), &jsonObjs2)
			totalEntriesCount += len(jsonObjs2.([]interface{}))
		}

		//Clean the collector, replace JSON Tag to get one big JSON array
		cleanCollector := bytes.Replace(collector, []byte("]["), []byte(","), -1)

		return cleanCollector, totalEntriesCount, err

		//If count of elements in first query is bigger or same than FilteredCount return
	} else {
		//Return
		return collector, totalEntriesCount, err
	}
}
