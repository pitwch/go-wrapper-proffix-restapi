package proffixrest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

//Error Struct PROFFIX REST-API
type ErrorStruct struct {
	Type    string `json:"Type"`
	Message string `json:"Message"`
	Fields  []struct {
		Reason  string `json:"Reason"`
		Name    string `json:"Name"`
		Message string `json:"Message"`
	} `json:"Fields"`
}

// GetFilteredCount returns the amount of total available entries in PROFFIX for this search query.
// Accepts a http.Header object
// Returns an integer
func GetFiltererCount(header http.Header) (total int) {
	type PxMetadata struct {
		FilteredCount int
	}
	var pxmetadata PxMetadata
	head := header.Get("pxmetadata")
	json.Unmarshal([]byte(head), &pxmetadata)

	return pxmetadata.FilteredCount
}

// errorFormatterPx does format the errors in a nicer way
func errorFormatterPx(ctx context.Context, c *Client, statuscode int, request io.Reader) (err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(request)
	errbyte := buf.Bytes()
	errstr := buf.String()
	log.Println(errstr)
	//No Autologout if deactivated in options
	if c.option.Autologout {
		//Error 404 is soft so no logout
		if statuscode != 404 {
			//Do Forced Logout...
			logoutstatus, err := c.Logout(ctx)
			logDebug(ctx, c, fmt.Sprintf("Logout after Error with Status: %v", logoutstatus))
			if err != nil {
				log.Printf("Error on Auto-Logout: %v", err)
			}
		}
	}

	//Define Error Struct
	parsedError := ErrorStruct{}

	//Try to parse JSON in ErrorStruct
	err = json.Unmarshal(errbyte, &parsedError)

	//If error on parse return plain text
	if err != nil {
		return fmt.Errorf("ERROR: %v", errstr)
	}
	if len(parsedError.Fields) > 0 {
		var errorFields []string
		for _, field := range parsedError.Fields {
			errorFields = append(errorFields, field.Name)
		}
		return fmt.Errorf("Status: %v, Type: %s, Message: %s, Fields: %v", statuscode, parsedError.Type, parsedError.Message, errorFields)

	}
	return fmt.Errorf("Status: %v, Type: %s, Message: %s", statuscode, parsedError.Type, parsedError.Message)
}

//WriteFile writes file from PROFFIX REST-API to local filepath
func WriteFile(filepath string, file io.ReadCloser) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write to file
	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}
	return nil
}
