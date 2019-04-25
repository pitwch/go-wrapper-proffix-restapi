package proffixrest

import (
	"context"
	"testing"
)

var address275 = Adresse{
	AdressNr: "275",
	Name:     "Test Adresse",
	Ort:      "Test",
	PLZ:      "8000",
}

var addressBatchTest = `[{
                   	"AdressNr": 276,
                   	"Name": "EYX AG",
                   	"Vorname": null,
                   	"Strasse": "Zollstrasse 2",
                   	"PLZ": "8000",
                   	"Ort": "Zürich"
                   }, {
                   	"Name": "Muster",
                   	"Vorname": "Hans",
                   	"Strasse": "Zürcherstrasse 2",
                   	"PLZ": "8000",
                   	"Ort": "Zürich"
                   }]`

func TestClient_SyncBatch(t *testing.T) {

	//New Context
	ctx := context.Background()

	//Connect
	pxrest, _ := ConnectTest(ctx, []string{"ADR"})

	//Delete Adress 274
	_, _, _, _ = pxrest.Delete(ctx, "ADR/Adresse/274")

	//Create Address 275
	_, _, _, _ = pxrest.Post(ctx, "ADR/Adresse", address275)

	result, _, _ := pxrest.SyncBatch(ctx, "ADR/Adresse", "AdressNr", []byte(addressBatchTest))

	var createdAdrNR string

	for _, res := range result {
		//Check Address 276
		if res[0] == "276" && res[1] != "204" {
			t.Errorf("AdressNr 276 should have PUT (204) Got '%v'", res[1])
		}
		//Check Address 276
		if res[0] != "276" && res[1] != "201" {
			t.Errorf("AdressNr %v should have POST (201) Got '%v'", res[0], res[1])
		}
		if res[0] != "276" && res[1] == "201" {
			createdAdrNR = res[0]
		}

	}
	//Delete created AdressNr
	_, _, statusDelete, errDeleted := pxrest.Delete(ctx, "ADR/Adresse/"+createdAdrNR)

	//Check HTTP Status of Logout. Should be 204
	if statusDelete != 204 {
		t.Errorf("Expected HTTP Status Code 204 for deleting AdressNr %v. Got '%v' with %v", createdAdrNR, statusDelete, errDeleted)
	}

	statuslogout, err := pxrest.Logout(ctx)

	//Check HTTP Status of Logout. Should be 204
	if statuslogout != 204 {
		t.Errorf("Expected HTTP Status Code 204. Got '%v'", err)
	}
}
