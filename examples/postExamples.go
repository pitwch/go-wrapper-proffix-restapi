package examples

import (
	"bytes"
	"context"
	"fmt"
	"log"
)

//var pxrest for Login already defined in getExamples.go

func postAdresseWithStruct() {

	//Create Context
	ctx := context.Background()

	//Connect to REST-API
	pxrest, err := Connect()

	//Define struct
	type Adresse struct {
		Name string
		Ort  string
		PLZ  string
	}
	var data = Adresse{"Muster GmbH", "Zürich", "8000"}

	//Query Endpoint ADR/Adresse with Headers
	rc, header, _, err := pxrest.Post(ctx, "ADR/Adresse", data)

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	fmt.Printf(resp, err)

	defer rc.Close()

	//Get Header Location
	fmt.Print(header.Get("Location"))
}

func postAdresseWithMap() {

	//Create Context
	ctx := context.Background()

	//Connect to REST-API
	pxrest, err := Connect()

	//Define map
	var data = map[string]interface{}{
		"Name":   "Muster GmbH",
		"Ort":    "Zürich",
		"Zürich": "8000",
	}

	//Query Endpoint ADR/Adresse with Headers
	rc, header, _, err := pxrest.Post(ctx, "ADR/Adresse", data)

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	fmt.Printf(resp, err)

	defer rc.Close()

	//Logout
	_, err = pxrest.Logout()

	//Log errors if there are
	if err != nil {
		log.Print(err)
	}

	//Get Header Location
	fmt.Print(header.Get("Location"))
}
