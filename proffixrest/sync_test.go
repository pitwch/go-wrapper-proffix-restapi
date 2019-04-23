package proffixrest

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"
)

type SampleArticle struct {
	ArtikelNr       string
	Bezeichnung1    string
	EinheitRechnung struct {
		EinheitNr string
	}
	EinheitLager struct {
		EinheitNr string
	}
}

func tokenGenerator() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

//Test Sync
func TestClient_Sync(t *testing.T) {

	ctx := context.Background()

	var artnrtest = tokenGenerator()
	//Connect
	pxrest, _ := ConnectTest(ctx)

	//Sample Articel
	art := SampleArticle{
		ArtikelNr:    artnrtest,
		Bezeichnung1: "Sample Article XYZTEST1",
		EinheitRechnung: struct {
			EinheitNr string
		}{"STK"},
		EinheitLager: struct {
			EinheitNr string
		}{"STK"},
	}

	_, _, status, err := pxrest.Sync(ctx, "LAG/Artikel", "LAG/Artikel", artnrtest, art)

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Connect. Got '%v'", err)
	}

	//Check status. Should be 201
	if status != 201 {
		t.Errorf("Expected 201. Got '%v'", status)
	}

	_, _, status, err = pxrest.Sync(ctx, "LAG/Artikel", "LAG/Artikel", artnrtest, art)

	//Check error. Should be nil
	if err != nil {
		t.Errorf("Expected no error for Connect. Got '%v'", err)
	}

	//Check status. Should be 204
	if status != 204 {
		t.Errorf("Expected 204. Got '%v'", status)
	}
}
