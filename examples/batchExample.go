package examples

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func getBatchExample() {

	//Define struct for Adresse
	type Adresse struct {
		AdressNr int
		Name     string
		Ort      string
		Plz      string
	}

	//Define Adresse as array (as we'll have many results)
	adressen := []Adresse{}

	//Connect to PROFFIX REST-API
	pxrest, err := Connect()
	if err != nil {
		fmt.Print(err)
	}

	//Set Params. As we just want some fields we define them on Fields param.
	params := url.Values{}
	params.Set("Fields", "AdressNr,Name,Ort,Plz")

	//Fire the batch request
	rc, total, _ := pxrest.GetBatch("ADR/Adresse", params, 35)

	if err != nil {
		fmt.Print(err)
	}

	//Unmarshal the result into our struct
	err = json.Unmarshal(rc, &adressen)
	if err != nil {
		fmt.Println("There was an error:", err)
	}

	//Print the result
	fmt.Print(adressen)

	//Check if we got them all...
	fmt.Printf("We got %v from %v Adressen!", len(adressen), total)

	//Logout
	_, err = pxrest.Logout()
}
