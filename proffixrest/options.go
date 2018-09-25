package proffixrest

import (
	"net/http"
	"time"
)

type Options struct {
	Key           string
	Version       string
	APIPrefix     string
	LoginEndpoint string
	UserAgent     string
	Timeout       time.Duration
	VerifySSL     bool
	Batchsize     int
	Log           bool
	Client        *http.Client
}
