package proffixrest

import (
	"bytes"
	"net/url"
	"testing"
)

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

func TestGetRequest(t *testing.T) {
	//Connect
	pxrest, err := ConnectTest()
	params := url.Values{}
	params.Set("Fields", "Adressnr")
	params.Set("Limit", "1")
	rc, headers, statuscode, err := pxrest.Get("ADR/Adresse", params)

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

	//Check if PxSessionId keeps the same
	_, testheader, _, _ := pxrest.Get("ADR/Adresse", url.Values{})

	pxsessionid1 := headers.Get("pxsessionid")
	pxsessionid2 := testheader.Get("pxsessionid")

	if pxsessionid1 != pxsessionid2 {
		t.Errorf("Session should be the same. Got different Session IDs.'%v' / '%v'", pxsessionid1, pxsessionid2)
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
		t.Errorf("Expected no error for Logout Request. Got '%v'", err)
	}

	//Check if response isn't empty
	if resp == "" {
		t.Errorf("Response shouldn't be empty. Got '%v'", resp)
	}

}
