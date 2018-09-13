package proffixrest

import (
	"bytes"
	"net/url"
	"path"
	"testing"
)

//Define struct
type Adresse struct {
	Name string
	Ort  string
	PLZ  string
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

func TestClient_NewClient(*testing.T) {

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
	var data = Adresse{"Muster GmbH", "Zürich", "8000"}

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

	//GET TESTs

	//Get Created AdressNr from Header, store in Var
	DemoAdressNr = path.Base(headers.Get("Location"))

	if DemoAdressNr == "" {
		t.Errorf("AdressNr should be in Location. Got '%v'", DemoAdressNr)
	}

	//Check if PxSessionId keeps the same
	_, testheader, _, _ := pxrest.Get("ADR/Adresse", url.Values{})

	pxsessionid1 := headers.Get("pxsessionid")
	pxsessionid2 := testheader.Get("pxsessionid")

	if pxsessionid1 != pxsessionid2 {
		t.Errorf("Session should be the same. Got different Session IDs.'%v' / '%v'", pxsessionid1, pxsessionid2)
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
	_, headers, statuscode, err = pxrest.Delete("ADR/Adresse/"+DemoAdressNr, url.Values{})

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
	statuslogout, err := pxrest.Logout("")

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Logout Request. Got '%v'", err)
	}

	//Check HTTP Status of Logout. Should be 204
	if statuslogout != 204 {
		t.Errorf("Expected HTTP Status Code 204. Got '%v'", err)
	}

}

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

func TestGetInfo(t *testing.T) {
	pxrest, err := ConnectTest()

	rc, err := pxrest.Info("")

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
