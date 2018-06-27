package examples

import (
	"fmt"
	"bytes"
	px "github.com/pitwch/go-wrapper-proffix-restapi/proffixrest"

	"net/url"
)

var pxrest, err = px.NewClient(
	"https://myserver.ch:999",
	"USR",
	"b62cce2fe18f7a156a9c719c57bebf0478a3d50f0d7bd18d9e8a40be2e663017",
	"DEMO",
	[]string{"ADR", "FIB"},
	&px.Options{},
)

func getInfoExample() {

	//User Helper Endpoint PRO/Info
	rc, err := pxrest.Info("112a5a90fe28b23ed2c776562a7d1043957b5b79fad242b10141254b4de59028")

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	fmt.Printf(resp, err)

	defer rc.Close()
}

func getDatabaseExample() {

	//User Helper Endpoint PRO/Datenbank
	rc, err := pxrest.Info("112a5a90fe28b23ed2c776562a7d1043957b5b79fad242b10141254b4de59028")

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	fmt.Printf(resp, err)

	defer rc.Close()
}

func getAdresseExample() {

	//Query endpoint ADR/Adresse/1
	rc, err := pxrest.Get("ADR/Adresse/1", url.Values{})

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	fmt.Printf(resp, err)

	defer rc.Close()
}

func getAdresseWithParamsExample() {

	//Set URL Params
	param := url.Values{}

	//Get all adresses with Vorname like Max
	param.Set("Filter", "Vorname@='Max'")

	//Query Endpoint ADR/Adresse
	rc, err := pxrest.Get("ADR/Adresse", param)

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	fmt.Printf(resp, err)

	defer rc.Close()
}
