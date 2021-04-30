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

// ReaderToString parses Reader to String
func ReaderToString(rc io.Reader) (str string, err error) {
	if rc != nil {
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(rc)
		str = buf.String()
	}
	return str, err
}

// ReaderToByte parses Reader to Bye
func ReaderToByte(rc io.Reader) (byt []byte, err error) {
	if rc != nil {
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(rc)
		byt = buf.Bytes()
	}
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

// GetMaps returns []map[string]interface{} from io.Reader
func GetMaps(rc io.Reader) (items []map[string]interface{}, err error) {
	resp, err := ReaderToByte(rc)
	err = json.Unmarshal(resp, &items)
	return items, err
}

// GetMap returns map[string]interface{} from io.Reader
func GetMap(rc io.Reader) (items map[string]interface{}, err error) {
	var it = map[string]interface{}{}
	res, err := ReaderToByte(rc)

	err = json.Unmarshal(res, &it)

	return it, err
}

// GetFileTokens returns map[string]string{} from io.Reader
func GetFileTokens(rc io.Reader, keyField string, fileField string) (files []map[string][]string, err error) {

	// Set Default
	if fileField == "" {
		fileField = "DateiNr"
	}

	// Iterate over items
	items, err := GetMaps(rc)
	for _, mp := range items {
		key := mp[keyField].(string)
		var tmpArr []string
		for k, v := range mp {
			if k == fileField {
				//var tmpMp map[string][]string
				tmpArr = append(tmpArr, v.(string))
			}
		}

		files = append(files, map[string][]string{key: tmpArr})
	}
	return files, err
}

// GetUsedLicences returns total amount of used licences
func GetUsedLicences(res io.Reader) (total float64, err error) {
	mp, err := GetMap(res)
	var Lizenzen = mp["Instanz"].(map[string]interface{})["Lizenzen"].([]interface{})

	total = 0
	for _, l := range Lizenzen {
		var liz = l.(map[string]interface{})["AnzahlInVerwendung"]
		total += liz.(float64)
	}
	return total, err
}
