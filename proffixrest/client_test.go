package proffixrest

import (
	"bytes"
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

//Store created Demo Address
var DemoAdressNr string

//Connect function
func ConnectTest() (pxrest *Client, err error) {

	//Use PROFFIX Demo Logins as Example
	pxrest, err = NewClient(
		"https://remote.proffix.net:11011/pxapi/v2",
		"Gast",
		"16ec7cb001be0525f9af1a96fd5ea26466b2e75ef3e96e881bcb7149cd7598da",
		"DEMODB",
		[]string{"ADR", "FIB"},
		&Options{
			Key: "16378f3e3bc8051435694595cbd222219d1ca7f9bddf649b9a0c819a77bb5e50"},
	)

	return pxrest, err
}

//Test Invalid Session
func TestClient_InvalidSession(t *testing.T) {
	//Connect
	pxrest, err := ConnectTest()

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Connect. Got '%v'", err)
	}

	//Set false Session ID
	Pxsessionid = "1234"

	//Set Params
	params := url.Values{}
	params.Set("Limit", "1")
	params.Set("Fields", "AdressNr")

	_, _, statuscode, err := pxrest.Get("ADR/Adresse", url.Values{})

	//Check status code; Should be 401
	if statuscode != 401 {
		t.Errorf("Expected HTTP Status Code 401. Got '%v'", statuscode)
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
	DemoAdressNr = path.Base(headers.Get("Location"))

	if DemoAdressNr == "" {
		t.Errorf("AdressNr should be in Location. Got '%v'", DemoAdressNr)
	}

	//PUT TESTs

	//Fill struct with Demo Data
	var update = Adresse{AdressNr: DemoAdressNr, Name: "Muster AG", Ort: "St. Gallen", PLZ: "9000"}

	//Put updated Data
	_, _, statuscode, err = pxrest.Put("ADR/Adresse/"+DemoAdressNr, update)

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
		t.Errorf("Expected no error for PUT Request. Got '%v'", err)
	}

	//GET TESTs

	//Check if PxSessionId keeps the same
	_, testheader, _, _ := pxrest.Get("ADR/Adresse", url.Values{})

	pxsessionid1 := headers.Get("pxsessionid")
	pxsessionid2 := testheader.Get("pxsessionid")

	if pxsessionid1 != pxsessionid2 {
		t.Errorf("Session should be the same. Got different Session IDs.'%v' / '%v'", pxsessionid1, pxsessionid2)
	}

	//Query Non Existing AdressNr
	_, _, statuscode, err = pxrest.Get("ADR/Adresse/123456789", url.Values{})

	//Check status code; Should be 404
	if statuscode != 404 {
		t.Errorf("Expected HTTP Status Code 404. Got '%v'", statuscode)
	}

	//Check error. Shouldn't be nil
	if err != nil {
		t.Errorf("Expected no errors for non existing GET Request. Got '%v'", err)
	}

	//Query Created AdressNr
	rc, headers, statuscode, err := pxrest.Get("ADR/Adresse/"+DemoAdressNr, url.Values{})

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
	buf := new(bytes.Buffer)
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
