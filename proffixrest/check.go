package proffixrest

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
)

// InfoStruct captures metadata returned by the PROFFIX REST-API info endpoint.
type InfoStruct struct {
	Version        string
	ServerZeit     string
	NeuesteVersion string
	Instanz        InstanzStruct
}

// InstanzStruct describes a PROFFIX instance and its licenses.
type InstanzStruct struct {
	InstanzNr string
	Name      string
	Lizenzen  []LizenzStruct
}

// LizenzStruct details a specific PROFFIX license entry.
type LizenzStruct struct {
	Name        string
	Bezeichnung string
	Anzahl      int
	Demo        bool
	Ablaufdatum string
}

// CheckAPI checks the PROFFIX REST API.
func (c *Client) CheckAPI(ctx context.Context, webservicepw string) (err error) {

	// Set timeout to 10
	c.option.Timeout = 10

	// Set Module ADR
	// c.Module = []string{"ADR"}

	// Check Login
	err = c.Login(ctx)
	if err != nil {
		return err
	}

	// Check if PxSessionID is available
	if c.GetPxSessionID() == "" {
		return &PxError{Message: "No PxSessionId found"}
	}

	closer, err := c.Info(ctx, webservicepw)
	if err != nil {
		return err
	}

	info := InfoStruct{}
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(closer)
	bytes := buf.Bytes()
	_ = json.Unmarshal(bytes, &info)

	for _, liz := range info.Instanz.Lizenzen {
		log.Println(liz.Name)
	}
	return nil

}
