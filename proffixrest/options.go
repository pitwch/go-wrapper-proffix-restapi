package proffixrest

import (
	"time"
)

type Options struct {
	Key           string
	Version       string
	APIPrefix     string
	LoginEndpoint string
	UserAgent     string
	Timeout       time.Duration
	VerifySSL       bool
}
