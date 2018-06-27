package examples

import (
	"bytes"
	"fmt"
)

//var pxrest for Login already defined in getExamples.go

func postAdresseWithStruct() {

	//Define struct
	type Adresse struct {
		Name string
		Ort  string
		PLZ  string
	}
	var data = Adresse{"Muster GmbH", "Zürich", "8000",
	}

	//Query Endpoint ADR/Adresse with Headers
	rc, header, err := pxrest.Post("ADR/Adresse", data)

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

	//Define map
	var data map[string]interface{} = map[string]interface{}{
		"Name":   "Muster GmbH",
		"Ort":    "Zürich",
		"Zürich": "8000",
	}

	//Query Endpoint ADR/Adresse with Headers
	rc, header, err := pxrest.Post("ADR/Adresse", data)

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	fmt.Printf(resp, err)

	defer rc.Close()


	//Get Header Location
	fmt.Print(header.Get("Location"))
}