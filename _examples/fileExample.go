package _examples

import (
	"context"
	"encoding/json"
	"github.com/pitwch/go-wrapper-proffix-restapi/proffixrest"
	"io/ioutil"
	"log"
)

func postfileExample() {

	// Create Context
	ctx := context.Background()

	// Connect to REST-API
	pxrest, err := Connect()

	// Get file
	fp := "C:/test_14.png"

	file, err := ioutil.ReadFile(fp)

	if err != nil {
		log.Print(err)
	}
	_, headers, _, err := pxrest.File(ctx, "test.png", file)

	// Convert location to ID
	id := proffixrest.ConvertLocationToID(headers)

	// Just use if theres an ID
	if id != "." {
		log.Printf("ID: %v", id)

		// Example Artikel
		mp := map[string]interface{}{}
		mp["DateiNr"] = id
		mp["Artikel"] = map[string]interface{}{"ArtikelNr": "TESTARTIKEL"}
		mp["Bezeichnung"] = "Test"

		jsn, err := json.Marshal(mp)
		if err != nil {
			log.Print(err)
		}
		log.Print(string(jsn))
		_, _, status, err := pxrest.Post(ctx, "LAG/Artikelbild", mp)

		log.Printf("Status: %v", status)

	}
	_, err = pxrest.Logout(ctx)
	if err != nil {
		log.Printf("Fehler bei Logout: %v", err)
	}
}
