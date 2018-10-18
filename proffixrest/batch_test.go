package proffixrest

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"
)

//Test GetBatch
func TestGetBatch(t *testing.T) {

	//New Context
	ctx := context.Background()

	//Define test struct for Adresse
	type Adresse struct {
		AdressNr int
		Name     string
		Ort      string
		Plz      string
		Land     string
	}

	//Define Adressen as array
	var adressen []Adresse

	//Connect
	pxrest, err := ConnectTest(ctx)

	//Set Params. As we just want some fields we define them on Fields param.
	params := url.Values{}
	params.Set("Fields", "AdressNr,Name,Ort,Plz")

	resultBatch, total, err := pxrest.GetBatch(ctx, "ADR/Adresse", params, 25)

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for GET Batch Request. Got '%v'", err)
	}

	//Check total. Should be greater than 0
	if total == 0 {
		t.Errorf("Expected some results. Got '%v'", total)
	}

	//Unmarshal the result into our struct
	err = json.Unmarshal(resultBatch, &adressen)

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Unmarshal GET Batch Request. Got '%v'", err)
	}

	//Logout
	statuslogout, err := pxrest.Logout(ctx)

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Logout. Got '%v'", err)
	}

	//Check HTTP Status of Logout. Should be 204
	if statuslogout != 204 {
		t.Errorf("Expected HTTP Status Code 204. Got '%v'", err)
	}

	//Compare received total vs. length of array. Should be equal.
	if !(len(adressen) == total) {
		t.Errorf("Total Results and Length of Array should be equal. Got '%v' vs. '%v'", len(adressen), total)

	}

	//Check default settings
	if pxrest.option.Batchsize != 200 {
		t.Errorf("Default batchsize should . Got '%v' vs. '%v'", len(adressen), total)
	}

}
