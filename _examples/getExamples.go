package _examples

import (
	"bytes"
	"context"
	"fmt"
	px "github.com/pitwch/go-wrapper-proffix-restapi/proffixrest"
	"log"
	"net/url"
)

//Put Client in Function
func Connect() (pxrest *px.Client, err error) {

	//Use PROFFIX Demo Logins as Example
	pxrest, err = px.NewClient(
		"https://remote.proffix.net:11011/pxapi/v2",
		"Gast",
		"16ec7cb001be0525f9af1a96fd5ea26466b2e75ef3e96e881bcb7149cd7598da",
		"DEMODB",
		[]string{"ADR", "FIB"},
		&px.Options{Timeout: 30},
	)

	return pxrest, err
}
func getInfoExample() {

	//Create Context
	ctx := context.Background()

	//Connect to REST-API
	pxrest, err := Connect()

	//User Helper Endpoint PRO/Info
	rc, err := pxrest.Info(ctx, "16378f3e3bc8051435694595cbd222219d1ca7f9bddf649b9a0c819a77bb5e50")

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	fmt.Printf(resp, err)

	defer rc.Close()

	//Logout
	_, err = pxrest.Logout(ctx)

	//Log errors if there are
	if err != nil {
		log.Print(err)
	}

}

func getDatabaseExample() {

	//Create Context
	ctx := context.Background()

	//Connect to REST-API
	pxrest, err := Connect()

	//User Helper Endpoint PRO/Datenbank
	rc, err := pxrest.Info(ctx, "16378f3e3bc8051435694595cbd222219d1ca7f9bddf649b9a0c819a77bb5e50")

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	fmt.Printf(resp, err)

	defer rc.Close()

	//Logout
	_, err = pxrest.Logout(ctx)

	//Log errors if there are
	if err != nil {
		log.Print(err)
	}
}

func getAdresseExample() {

	//Create Context
	ctx := context.Background()

	//Connect to REST-API
	pxrest, err := Connect()

	//Query endpoint ADR/Adresse/1 without Header as not needed
	rc, _, _, err := pxrest.Get(ctx, "ADR/Adresse/1", url.Values{})

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	fmt.Printf(resp, err)

	defer rc.Close()

	//Logout
	_, err = pxrest.Logout(ctx)

	//Log errors if there are
	if err != nil {
		log.Print(err)
	}

}

func getAdresseWithParamsExample() {

	//Create Context
	ctx := context.Background()

	//Connect to REST-API
	pxrest, err := Connect()

	//Set URL Params
	param := url.Values{}

	//Get all adresses with Vorname like Max
	param.Set("Filter", "Vorname@='Max'")

	//Query Endpoint ADR/Adresse
	rc, _, _, err := pxrest.Get(ctx, "ADR/Adresse", param)

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	fmt.Printf(resp, err)

	defer rc.Close()

	//Logout
	_, err = pxrest.Logout(ctx)

	//Log errors if there are
	if err != nil {
		log.Print(err)
	}
}
