package proffixrest

import (
	"log"
	"net/http"
	"time"
)

type Options struct {
	Key           string        // API-Key for PROFFIX REST-API
	Version       string        // Version of API to use. Default is v3
	APIPrefix     string        // API Prefix. Default is /pxapi/
	LoginEndpoint string        // Login Endpoint. Default is PRO/Login
	UserAgent     string        // User Agent. Default is go-wrapper-proffix-restapi
	Timeout       time.Duration // Timeout. Default is 15 seconds
	VerifySSL     bool          // Verifies SSL Cert of REST-API. Default is true
	Batchsize     int
	Log           bool
	Autologout    bool
	VolumeLicence bool         // If API should use Volume Licencing
	HTTPClient    *http.Client // Optional custom HTTP client to use
	Logger        *log.Logger  // Optional logger; overrides Log flag when provided
}
