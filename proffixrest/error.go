package proffixrest

import (
	"encoding/json"
	"fmt"
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
		var b strings.Builder
		for i, f := range e.Fields {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(f.Name)
		}
		return e.Message + " (" + b.String() + ")"
	}
	return e.Message
}

// NewPxError creates a new PxError
func NewPxError(rc io.Reader, status int, endpoint string) *PxError {

	// Initialize with known context so we can craft a helpful message if body is empty
	pxerr := PxError{Status: status, Endpoint: endpoint}

	if rc != nil {
		dec := json.NewDecoder(rc)
		_ = dec.Decode(&pxerr)
		// Ensure status/endpoint are preserved even if JSON doesn't contain them
		pxerr.Status = status
		pxerr.Endpoint = endpoint
	}

	// If API did not provide a message, set a sensible default to avoid empty error strings
	if pxerr.Message == "" {
		switch {
		case pxerr.Status > 0:
			pxerr.Message = fmt.Sprintf("HTTP %d from %s", pxerr.Status, pxerr.Endpoint)
		case pxerr.Endpoint != "":
			pxerr.Message = fmt.Sprintf("request to %s failed", pxerr.Endpoint)
		default:
			pxerr.Message = "unknown error"
		}
	}

	return &pxerr
}
