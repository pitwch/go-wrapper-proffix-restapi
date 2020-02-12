package proffixrest

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
)

type InfoStruct struct {
	Version        string
	ServerZeit     string
	NeuesteVersion string
	Instanz        InstanzStruct
}

type InstanzStruct struct {
	InstanzNr string
	Name      string
	Lizenzen  []LizenzStruct
}

type LizenzStruct struct {
	Name               string
	Bezeichnung        string
	Anzahl             int
	AnzahlInVerwendung int
	Demo               bool
	Ablaufdatum        string
}

// CheckAPI checks the PROFFIX REST API
// Accepts client
// Returns err
func (c *Client) CheckApi(ctx context.Context, webservicepw string) (err error) {

	//Set timeout to 10
	c.option.Timeout = 10

	//Set Module ADR
	c.Module = []string{"ADR"}

	//Check Login
	err = c.Login(ctx)
	if err != nil {
		return err
	}

	//Check if PxSessionId is available
	if c.GetPxSessionId() == "" {
		return &PxError{Message: "No PxSessionId found"}
	}

	closer, err := c.Info(ctx, webservicepw)
	if err != nil {
		return err
	}

	info := InfoStruct{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(closer)
	bytes := buf.Bytes()
	json.Unmarshal(bytes, &info)

	for _, liz := range info.Instanz.Lizenzen {
		log.Println(liz.Name)
	}
	return nil

}
