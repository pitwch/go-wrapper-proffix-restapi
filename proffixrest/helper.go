package proffixrest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// GetFilteredCount returns the total available entries reported by PROFFIX for a search query.
func GetFilteredCount(header http.Header) (total int) {
	type PxMetadata struct {
		FilteredCount int
	}
	var pxmetadata PxMetadata
	head := header.Get("pxmetadata")
	_ = json.Unmarshal([]byte(head), &pxmetadata) // Ignore error, return 0 on failure

	return pxmetadata.FilteredCount
}

// ReaderToString parses an io.Reader into its string representation.
func ReaderToString(rc io.Reader) (str string, err error) {
	if rc != nil {
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(rc)
		str = buf.String()
	}
	return str, err
}

// ReaderToByte parses an io.Reader into a byte slice.
func ReaderToByte(rc io.Reader) (byt []byte, err error) {
	if rc != nil {
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(rc)
		byt = buf.Bytes()
	}
	return byt, err
}

// WriteFile writes file from PROFFIX REST-API to local filepath.
func WriteFile(filePath string, file io.ReadCloser) (err error) {
	cleanPath := filepath.Clean(filePath)

	// Create the file
	out, err := os.Create(cleanPath) // #nosec G304 -- caller controls destination path by design
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	// Write to file
	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}
	return nil
}

// GetMaps returns []map[string]interface{} from io.Reader.
func GetMaps(rc io.Reader) (items []map[string]interface{}, err error) {
	resp, err := ReaderToByte(rc)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(resp, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// GetMap returns map[string]interface{} from io.Reader.
func GetMap(rc io.Reader) (items map[string]interface{}, err error) {
	var it = map[string]interface{}{}
	res, err := ReaderToByte(rc)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(res, &it); err != nil {
		return nil, err
	}
	return it, nil
}

// GetFileTokens returns map[string]string{} from io.Reader.
func GetFileTokens(rc io.Reader, keyField string, fileField string) (files []map[string][]string, err error) {
	if fileField == "" {
		fileField = "DateiNr"
	}

	items, err := GetMaps(rc)
	if err != nil {
		return nil, err
	}

	for _, mp := range items {
		key := mp[keyField].(string)
		var tmpArr []string
		for k, v := range mp {
			if k == fileField {
				tmpArr = append(tmpArr, v.(string))
			}
		}

		files = append(files, map[string][]string{key: tmpArr})
	}
	return files, nil
}

// GetUsedLicences returns total amount of used licences.
func GetUsedLicences(res io.Reader) (total float64, err error) {
	mp, err := GetMap(res)
	if err != nil {
		return 0, err
	}
	lizenzen := mp["Instanz"].(map[string]interface{})["Lizenzen"].([]interface{})

	for _, l := range lizenzen {
		liz := l.(map[string]interface{})["AnzahlInVerwendung"]
		total += liz.(float64)
	}
	return total, nil
}
