package examples

import (
	"bytes"
	"context"
	"fmt"
	px "github.com/pitwch/go-wrapper-proffix-restapi/proffixrest"
	"google.golang.org/appengine/urlfetch"
	"log"
	"net/url"
)

//Put Client in Function
func ConnectAppEngine() (pxrest *px.Client, ctx context.Context, err error) {

	//Set Context
	ctx = context.Background()

	//Use URLFetch as HTTP Client as this is mandatory on Google App Engine Standard
	px.DefaultHTTPClient = urlfetch.Client(ctx)

	//Use PROFFIX Demo Logins as Example
	pxrest, err = px.NewClient(
		"https://remote.proffix.net:11011/pxapi/v2",
		"Gast",
		"16ec7cb001be0525f9af1a96fd5ea26466b2e75ef3e96e881bcb7149cd7598da",
		"DEMODB",
		[]string{"ADR", "FIB"},
		&px.Options{Timeout: 30, VerifySSL: false},
	)

	return pxrest, ctx, err
}

func getAdresseGoogleAppEngineExample() {

	//Connect to REST-API
	pxrest, ctx, err := ConnectAppEngine()

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
