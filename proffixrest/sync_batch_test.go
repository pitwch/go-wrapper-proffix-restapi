package proffixrest

import (
	"context"
	"encoding/json"
	"testing"
)

var artikel99999 = `{
	"ArtikelNr": "99999",
	"Bezeichnung1": "Testartikel 99999",
	"Bezeichnung2": "Extras, Test, Spezial",
	"Verkaufspreis1": "145.00",
	"Steuercode1": {},
	"Gewicht": "0.5",
	"EinheitLager": {
		"EinheitNr": "STK"
	},
	"EinheitRechnung": {
		"EinheitNr": "STK"
	},
	"Waehrung": {
		"WaehrungNr": "CHF"
	},
	"Ertragskonto": {},
	"Steuercode": {},
	"KeinBestand": true
}`

var artikelBatchTest = `[{
	"ArtikelNr": "99999",
	"Bezeichnung1": "Testartikel 99999 updated",
	"Bezeichnung2": "Telefonsystem: AS-Serie,Telefonkategorie: Digital,8 Up0-Ports für AGFEO Up0-Systemtelefone,",
	"Verkaufspreis1": "112.00",
	"Steuercode1": {},
	"Gewicht": "0.25",
	"EinheitLager": {
		"EinheitNr": "STK"
	},
	"EinheitRechnung": {
		"EinheitNr": "STK"
	},
	"Waehrung": {
		"WaehrungNr": "CHF"
	},
	"Ertragskonto": {},
	"Steuercode": {},
	"KeinBestand": true
}, {
	"ArtikelNr": "66666",
	"Bezeichnung1": "Testartikel 66666",
	"Bezeichnung2": "Telefonsystem: AS-Systeme,Montage: Tisch; Wand,1x S0 extern,3x S0 schaltbar (intern/extern),",
	"Verkaufspreis1": "256.00",
	"Steuercode1": {},
	"Gewicht": "1.98",
	"EinheitLager": {
		"EinheitNr": "STK"
	},
	"EinheitRechnung": {
		"EinheitNr": "STK"
	},
	"Waehrung": {
		"WaehrungNr": "CHF"
	},
	"Ertragskonto": {},
	"Steuercode": {},
	"KeinBestand": true
}]`

func TestClient_SyncBatch(t *testing.T) {

	//New Context
	ctx := context.Background()

	//Connect
	pxrest, _ := ConnectTest(ctx, []string{"LAG"})

	//Delete Artikel 99999
	_, _, _, _ = pxrest.Delete(ctx, "LAG/Artikel/99999")

	//Create Artikel 99999
	artikelMap := make(map[string]interface{})

	err := json.Unmarshal([]byte(artikel99999), &artikelMap)
	_, _, _, _ = pxrest.Post(ctx, "LAG/Artikel", artikelMap)

	created, updated, failed, errors, total, err := pxrest.SyncBatch(ctx, "LAG/Artikel", "ArtikelNr", false, []byte(artikelBatchTest))

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
	if len(created) > 0 {
		//Delete created AdressNr
		_, _, statusDelete, errDeleted := pxrest.Delete(ctx, "LAG/Artikel/"+created[0])

		//Check HTTP Status of Logout. Should be 204
		if statusDelete != 204 {
			t.Errorf("Expected HTTP Status Code 204 for deleting AdressNr %v. Got '%v' with %v", created[0], statusDelete, errDeleted)
		}
	}
	statuslogout, err := pxrest.Logout(ctx)

	//Check HTTP Status of Logout. Should be 204
	if statuslogout != 204 {
		t.Errorf("Expected HTTP Status Code 204. Got '%v'", err)
	}
}
