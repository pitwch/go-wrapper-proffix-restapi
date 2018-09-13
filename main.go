package main

import (
	"fmt"
	px "github.com/pitwch/go-wrapper-proffix-restapi/proffixrest"
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
	//Buffer decode for plain text response

	//Logout
	_, err = pxrest.Logout()
}
