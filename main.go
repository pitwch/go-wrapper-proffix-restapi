package main

import (
	"bytes"
	"fmt"
	px "github.com/pitwch/go-wrapper-proffix-restapi/proffixrest"
	"net/url"
)

func Connect() (pxrest *px.Client, err error) {

	//Use PROFFIX Demo Logins as Example
	pxrest, err = px.NewClient(
		"https://remote.proffix.net:11011/pxapi/v2",
		"Gast",
		"16ec7cb001be0525f9af1a96fd5ea26466b2e75ef3e96e881bcb7149cd7598da",
		"DEMODB",
		[]string{"ADR"},
		&px.Options{},
	)

	return pxrest, err
}

func main() {

	pxrest, err := Connect()
	if err != nil {
		fmt.Print(err)
	}
	rc, _, status, err := pxrest.Get("ADR/Adresse", url.Values{})
	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	fmt.Printf(resp, err)
	defer rc.Close()
	fmt.Print(status)

	fmt.Print(resp)
	//Logout
	_, err = pxrest.Logout("")
}
