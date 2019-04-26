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

	created, updated, failed, errors, total, err := pxrest.SyncBatch(ctx, "ADR/Adresse", "AdressNr", true, []byte(addressBatchTest))

	//Calculate created / updated
	checkCreated := total - len(updated)
	checkUpdated := total - len(created)

	//Check Created
	if len(created) != checkCreated {
		t.Errorf("Created should be %v. Is %v", len(created), checkCreated)
	}

	//Check Updated
	if len(updated) != checkUpdated {
		t.Errorf("Updated should be %v. Is %v", len(updated), checkUpdated)
	}

	//Check Failed
	if len(failed) != 0 {
		t.Errorf("Failed should be 0. Is %v : %v", len(failed), failed)
	}

	//Check Errors
	if len(errors) != 0 {
		t.Errorf("Errors should be 0. Is %v : %v", len(errors), errors)
	}
	//Delete created AdressNr
	_, _, statusDelete, errDeleted := pxrest.Delete(ctx, "ADR/Adresse/"+created[0])

	//Check HTTP Status of Logout. Should be 204
	if statusDelete != 204 {
		t.Errorf("Expected HTTP Status Code 204 for deleting AdressNr %v. Got '%v' with %v", created[0], statusDelete, errDeleted)
	}

	statuslogout, err := pxrest.Logout(ctx)

	//Check HTTP Status of Logout. Should be 204
	if statuslogout != 204 {
		t.Errorf("Expected HTTP Status Code 204. Got '%v'", err)
	}
}
