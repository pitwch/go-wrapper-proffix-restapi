package proffixrest

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/spf13/cast"
	"net/url"
	"path"
	"testing"
)

//Define struct (AdressNr as string just for convenience)
type Adresse struct {
	AdressNr string `json:"AdressNr,omitempty"`
	Name     string `json:"Name,omitempty"`
	Ort      string `json:"Ort,omitempty"`
	PLZ      string `json:"PLZ,omitempty"`
}

//Connect function
func ConnectTest() (pxrest *Client, err error) {

	//New Context
	ctx := context.Background()

	//Use PROFFIX Public Demo Login as Example
	pxrest, err = NewClient(ctx,
		"https://remote.proffix.net:11011/pxapi/v2",
		"Gast",
		"16ec7cb001be0525f9af1a96fd5ea26466b2e75ef3e96e881bcb7149cd7598da",
		"DEMODB",
		[]string{"ADR", "FIB"},
		&Options{
			Key:       "16378f3e3bc8051435694595cbd222219d1ca7f9bddf649b9a0c819a77bb5e50",
			VerifySSL: false},
	)
	return pxrest, err
}

//Helper function for checking if is type Adresse
func isAdressArray(t interface{}) bool {
	switch t.(type) {
	case []Adresse:
		return true
	default:
		return false
	}
}

//Helper function for checking if is type Options
func isOption(t interface{}) bool {
	switch t.(type) {
	case Options:
		return true
	default:
		return false
	}
}

//Helper function for checking if is type Client
func isClient(t interface{}) bool {
	switch t.(type) {
	case Client:
		return true
	default:
		return false
	}
}

//Test all Requests in one Session
func TestClient_Requests(t *testing.T) {

	//POST TESTs

	//Fill struct with Demo Data
	var data = Adresse{Name: "Muster GmbH", Ort: "Zürich", PLZ: "8000"}

	//Connect
	pxrest, err := ConnectTest()
	_, headers, statuscode, err := pxrest.Post("ADR/Adresse", data)

	//Check status code; Should be 201
	if statuscode != 201 {
		t.Errorf("Expected HTTP Status Code 201. Got '%v'", statuscode)
	}

	//Check PXSessionId; Shouldn't be empty
	if headers.Get("pxsessionid") == "" {
		t.Errorf("Expected PxSessionId in Header. Not found: '%v'", headers.Get("pxsessionid"))
	}

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for GET Request. Got '%v'", err)
	}

	//Get Created AdressNr from Header, store in Var
	DemoAdressNr := path.Base(headers.Get("location"))

	if DemoAdressNr == "" {
		t.Errorf("AdressNr should be in Location. Got '%v'", DemoAdressNr)
	}

	//PUT TESTs

	//Fill struct with Demo Data
	var update = Adresse{AdressNr: DemoAdressNr, Name: "Muster AG", Ort: "St. Gallen", PLZ: "9000"}

	//Put updated Data
	rc, _, statuscode, err := pxrest.Put("ADR/Adresse/"+cast.ToString(DemoAdressNr), update)

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	rupd := buf.String()

	//Check status code; Should be 204
	if statuscode != 204 {
		t.Errorf("Expected HTTP Status Code 204. Got '%v'. Message: '%v'", statuscode, rupd)
		if statuscode == 500 {
			t.Errorf("Endpoint: 'ADR/Adresse/%v', Data: '%v'", DemoAdressNr, update)
		}
	}

	//Check PXSessionId; Shouldn't be empty
	if headers.Get("pxsessionid") == "" {
		t.Errorf("Expected PxSessionId in Header. Not found: '%v'.  Message: '%v'", headers.Get("pxsessionid"), rupd)
	}

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for PUT Request. Got '%v'. Message: '%v'", err, rupd)
	}

	//GET TESTs

	//Query Non Existing AdressNr
	_, _, statuscode, err = pxrest.Get("ADR/Adresse/123456789", url.Values{})

	//Check status code; Should be 404
	if statuscode != 404 {
		t.Errorf("Expected HTTP Status Code 404. Got '%v'", statuscode)
	}

	//Check error. Shouldn't be nil
	if err == nil {
		t.Errorf("Expected errors for non existing GET Request. Got '%v'", err)
	}

	//Query Created AdressNr
	rc, headers, statuscode, err = pxrest.Get("ADR/Adresse/"+DemoAdressNr, url.Values{})

	//Check status code; Should be 200
	if statuscode != 200 {
		t.Errorf("Expected HTTP Status Code 200. Got '%v'", statuscode)
	}

	//Check PXSessionId; Shouldn't be empty
	if headers.Get("pxsessionid") == "" {
		t.Errorf("Expected PxSessionId in Header. Not found: '%v'", headers.Get("pxsessionid"))
	}

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for GET Request. Got '%v'", err)
	}

	//Check response. Shouldn't be empty

	//Buffer decode for plain text response
	buf = new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	if resp == "" {
		t.Errorf("Response shouldn't be empty. Got '%v'", resp)
	}

	//DELETE TESTs

	//Delete the created Address
	_, headers, statuscode, err = pxrest.Delete("ADR/Adresse/" + DemoAdressNr)

	//Check status code; Should be 204
	if statuscode != 204 {
		t.Errorf("Expected HTTP Status Code 204. Got '%v'", statuscode)
	}

	//Check PXSessionId; Shouldn't be empty
	if headers.Get("pxsessionid") == "" {
		t.Errorf("Expected PxSessionId in Header. Not found: '%v'", headers.Get("pxsessionid"))
	}

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for DELETE Request. Got '%v'", err)
	}

	//Check Logout
	statuslogout, err := pxrest.Logout()

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Logout Request. Got '%v'", err)
	}

	//Check HTTP Status of Logout. Should be 204
	if statuslogout != 204 {
		t.Errorf("Expected HTTP Status Code 204. Got '%v'", err)
	}

}

//Test all Requests in one Session
func TestClient_AdvancedFilters(t *testing.T) {

	//Define test struct for Adresse
	type Adresse struct {
		AdressNr int
		Name     string
		Ort      string
		Plz      string
		Land     struct {
			LandNr      string
			Bezeichnung string
		}
	}

	//Define Adressen as array
	var adressen []Adresse

	//Connect
	pxrest, err := ConnectTest()

	//Set Advanced Params.
	params := url.Values{}
	params.Set("fields", "AdressNr,Name,Ort,Plz,Land")
	params.Set("filter", "Name@='a',Land=='CH',Ort!='Zürich'")
	params.Set("sort", "Plz,Name,AdressNr")

	//Fire Request
	resultFilters, _, statuscode, err := pxrest.Get("ADR/Adresse", params)

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for GET Batch Request. Got '%v'", err)
	}

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(resultFilters)
	respFilter := buf.Bytes()
	stringFilter := buf.String()
	//Unmarshal the result into our struct
	err = json.Unmarshal(respFilter, &adressen)

	//Check status code; Should be 200
	if statuscode != 200 {
		t.Errorf("Expected HTTP Status Code 200. Got '%v'. Message: %v", statuscode, stringFilter)
	}

	//Check if results match to filter Ort
	if adressen[0].Ort == "Zürich" {
		t.Errorf("Filter 'Ort' doesnt work. Got '%v' - shouldnt be.", adressen[0].Ort)
	}

	//Check if results match to filter Land
	if adressen[0].Land.LandNr != "CH" {
		t.Errorf("Filter 'Land' doesnt work. Got '%v' - shouldnt be.", adressen[0].Land)
	}

	//Check if order works
	if adressen[0].AdressNr > adressen[1].AdressNr {
		t.Errorf("Order doesnt work. Adressnr '%v' is > than Adressnr '%v'.", adressen[0].AdressNr, adressen[1].AdressNr)
	}

	//Logout
	pxrest.Logout()

}

//Test Get Database with Key from Options
func TestGetDatabase(t *testing.T) {
	pxrest, err := ConnectTest()

	rc, err := pxrest.Database("")

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Database Request. Got '%v'", err)
	}

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Database Request. Got '%v'", err)
	}

	//Check if response isn't empty
	if resp == "" {
		t.Errorf("Response shouldn't be empty. Got '%v'", resp)
	}

	//Test with key instead value from options
	_, err = pxrest.Database("16378f3e3bc8051435694595cbd222219d1ca7f9bddf649b9a0c819a77bb5e50")

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Database Request with key. Got '%v'", err)
	}
}

//Test Get Info with Key set in Function
func TestGetInfo(t *testing.T) {
	pxrest, err := ConnectTest()

	rc, err := pxrest.Info("16378f3e3bc8051435694595cbd222219d1ca7f9bddf649b9a0c819a77bb5e50")

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Info Request. Got '%v'", err)
	}

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Info Request. Got '%v'", err)
	}

	//Check if response isn't empty
	if resp == "" {
		t.Errorf("Response shouldn't be empty. Got '%v'", resp)
	}

}

//Test GetBatch
func TestGetBatch(t *testing.T) {

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
	pxrest, err := ConnectTest()

	//Set Params. As we just want some fields we define them on Fields param.
	params := url.Values{}
	params.Set("Fields", "AdressNr,Name,Ort,Plz")

	resultBatch, total, err := pxrest.GetBatch("ADR/Adresse", params, 25)

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
	statuslogout, err := pxrest.Logout()

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

//Check structs
func TestStructs(t *testing.T) {

	pxrest, _ := ConnectTest()

	//Set test vars
	q := *pxrest
	w := *pxrest.option

	//Check if pxrest is Options
	if isOption(w) == false {
		t.Errorf("w not of type Option")
	}

	//Check if pxrest.option is Options
	if isClient(q) == false {
		t.Errorf("q is not of type Client")
	}

	pxrest.Logout()
}

func TestClientError(t *testing.T) {
	urlstring, _ := url.Parse("https://bing.com")
	client := Client{
		restURL:   urlstring,
		Benutzer:  "test",
		Passwort:  "1234",
		Datenbank: "demo",
		Module:    nil,
		option:    &Options{LoginEndpoint: "Test"},
		rawClient: nil,
	}

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Client should have panicked!")
			}
		}()
		// This function should cause a panic
		_, _, _, _ = client.Get("test", nil)
	}()

}

func TestClient_LoginWithFalsPxSessionId(t *testing.T) {

	pxrest, err := ConnectTest()

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Logout. Got '%v'", err)
	}

	pxrest.Logout()
}
