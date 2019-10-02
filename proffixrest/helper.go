package proffixrest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// GetFilteredCount returns the amount of total available entries in PROFFIX for this search query.
// Accepts a http.Header object
// Returns an integer
func GetFiltererCount(header http.Header) (total int) {
	type PxMetadata struct {
		FilteredCount int
	}
	var pxmetadata PxMetadata
	head := header.Get("pxmetadata")
	_ = json.Unmarshal([]byte(head), &pxmetadata)

	return pxmetadata.FilteredCount
}

func ReaderToString(rc io.Reader) (str string, err error) {
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(rc)
	str = buf.String()
	return str, err
}

func ReaderToByte(rc io.Reader) (byt []byte, err error) {
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(rc)
	byt = buf.Bytes()
	return byt, err
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
