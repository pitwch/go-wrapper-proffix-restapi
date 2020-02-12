package proffixrest

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

// PxError defines the Error struct of default Error of PROFFIX REST-API
type PxError struct {
	Endpoint string           `json:"Endpoint"`
	Status   int              `json:"Status"`
	Type     string           `json:"Type"`
	Message  string           `json:"Message"`
	Fields   []PxInvalidField `json:"Fields"`
}

// PxInvalidField defines the Error struct of Fields Error of PROFFIX REST-API
type PxInvalidField struct {
	Reason  string `json:"Reason"`
	Name    string `json:"Name"`
	Message string `json:"Message"`
}

// Check if PxError is Null
func (e *PxError) isNull() bool {
	return e.Message == ""
}

// Check if is Type "INVALID_FIELDS"
func (e *PxError) isInvalidFields() bool {
	return e.Type == "INVALID_FIELDS"
}

// Check if is Type "NOT_FOUND"
func (e *PxError) isNotFound() bool {
	return e.Type == "NOT_FOUND"
}

// Formats error default
func (e *PxError) Error() string {
	if len(e.Fields) > 0 {
		var errFields = []string{}
		for _, f := range e.Fields {
			errFields = append(errFields, f.Name)
		}
		return e.Message + " (" + strings.Join(errFields, ",") + ")"
	} else {
		return e.Message
	}
}

// NewPxError creates a new PxError
func NewPxError(rc io.Reader, status int, endpoint string) *PxError {

	pxerr := PxError{}

	if rc != nil {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(rc)
		rcByte := buf.Bytes()
		_ = json.Unmarshal(rcByte, &pxerr)
	}
	pxerr.Status = status
	pxerr.Endpoint = endpoint

	if pxerr.isNull() {
		return &PxError{}
	} else {
		return &pxerr
	}
}
