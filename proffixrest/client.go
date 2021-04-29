package proffixrest

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

//Version of Wrapper
const (
	Version = "1.7.5"
)

// DefaultHTTPTransport is an http.RoundTripper that has DisableKeepAlives set true.
var DefaultHTTPTransport = &http.Transport{
	DisableKeepAlives: true,
}

// DefaultHTTPClient is an http.Client with the DefaultHTTPTransport and (Cookie) Jar set nil.
var DefaultHTTPClient = &http.Client{
	Jar:       nil,
	Transport: DefaultHTTPTransport,
}

// PxSessionId contains the latests PxSessionId from PROFFIX REST-API
var PxSessionId string

var isLoggedIn bool

//Client Struct
type Client struct {
	restURL   *url.URL
	Benutzer  string
	Passwort  string
	Datenbank string
	Module    []string
	option    *Options
	client    *http.Client
}

//Login Struct
type LoginStruct struct {
	Benutzer  string         `json:"Benutzer"`
	Passwort  string         `json:"Passwort"`
	Datenbank DatabaseStruct `json:"Datenbank"`
	Module    []string       `json:"Module,omitempty"`
}

//Database Struct
type DatabaseStruct struct {
	Name string `json:"Name"`
}

// Building new Client
func NewClient(RestURL string, apiUser string, apiPassword string, apiDatabase string, apiModule []string, options *Options) (*Client, error) {
	restURL, err := url.Parse(RestURL)
	if err != nil {
		return nil, &PxError{Message: fmt.Sprintf("%v", err)}
	}

	if options == nil {
		options = &Options{}
	}

	// Default to v3
	if options.Version == "" {
		options.Version = "v3"
	}

	if options.APIPrefix == "" {
		options.APIPrefix = "/pxapi/"
	}

	if options.LoginEndpoint == "" {
		options.LoginEndpoint = "PRO/Login"
	}

	if options.UserAgent == "" {
		options.UserAgent = "go-wrapper-proffix-restapi " + Version
	}

	// If VolumeLicence enabled set Modules to "VOL"
	if options.VolumeLicence == true {
		apiModule = []string{"VOL"}
	}

	// Set default batchsize for batch requests
	if options.Batchsize == 0 {
		options.Batchsize = 200
	}

	// Disable Cert Verfication
	if options.VerifySSL == false {
		DefaultHTTPTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if DefaultHTTPClient == nil {

		DefaultHTTPClient = http.DefaultClient

		// Set options for Timeout
		DefaultHTTPClient.Timeout = options.Timeout

	}

	path := options.APIPrefix + options.Version + "/"
	restURL.Path = path

	return &Client{
		restURL:   restURL,
		Benutzer:  apiUser,
		Passwort:  apiPassword,
		Datenbank: apiDatabase,
		Module:    apiModule,
		option:    options,
		client:    DefaultHTTPClient,
	}, nil
}

// Function for creating new PxSessionId
func (c *Client) createNewPxSessionId(ctx context.Context) (sessionid string, err error) {

	// Check URL, else exit
	_, err = url.ParseRequestURI(c.restURL.String())
	if err != nil {
		return "", &PxError{Message: "URL in wrong format"}
	}

	// Build Login URL
	urlstr := c.restURL.String() + c.option.LoginEndpoint

	// Create Login Body
	data := LoginStruct{c.Benutzer, c.Passwort, DatabaseStruct{c.Datenbank}, c.Module}

	// Encode Login Body to JSON
	body := new(bytes.Buffer)
	encoder := json.NewEncoder(body)
	if err := encoder.Encode(data); err != nil {
		return "", &PxError{Message: fmt.Sprintf("%v", err)}
	}

	// Build Login Request
	req, err := http.NewRequest("POST", urlstr, body)
	if err != nil {
		return "", &PxError{Message: "Error on builiding Request"}
	}

	// Set Login Header
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		if resp == nil {
			return "", NewPxError(nil, 0, c.option.LoginEndpoint)
		} else {
			return "", NewPxError(resp.Body, resp.StatusCode, c.option.LoginEndpoint)

		}

	} else if resp.StatusCode != 201 {
		return "", NewPxError(resp.Body, resp.StatusCode, c.option.LoginEndpoint)

	} else {
		// Set isLoggedIn Status
		isLoggedIn = true
		// Return pxsessionid
		return resp.Header.Get("pxsessionid"), &PxError{}
	}
}

// Function for Login
// Helper Function for Login to PROFFIX REST-API
// Login is automatically done
func (c *Client) Login(ctx context.Context) error {

	// If Pxsessionid doesnt yet exists create a new one
	if isLoggedIn == false {
		sessionid, err := c.createNewPxSessionId(ctx)

		PxSessionId = sessionid
		return err
		// If Pxsessionid already exists return stored value
	} else {
		return &PxError{}
	}
}

// updatePxSessionId updates the stored PxSessionId
func (c *Client) updatePxSessionId(header http.Header) {

	// Just update if PxSessionId in Header is not empty
	if header.Get("pxsessionid") != "" {
		PxSessionId = header.Get("pxsessionid")
	}

}

// LOGOUT Request for PROFFIX REST-API
// Accepts PxSession ID as Input or if left empty uses SessionId from active Session
// Returns Statuscode,error
func (c *Client) Logout(ctx context.Context) (int, error) {
	//Just logout if we have a valid PxSessionid
	if PxSessionId != "" {

		//Delete Login Object from PROFFIX REST-API
		req, _, statuscode, _ := c.request(ctx, "DELETE", c.option.LoginEndpoint, url.Values{}, false, nil)

		//Set Pxsessionid to empty string
		PxSessionId = ""

		if statuscode == 204 {
			isLoggedIn = false
			return statuscode, nil

		} else {
			return statuscode, NewPxError(req, statuscode, c.option.LoginEndpoint)

		}
	}
	return 0, nil

}

// Request Method
// Building the Request Method for Client
func (c *Client) request(ctx context.Context, method, endpoint string, params url.Values, isFile bool, data interface{}) (io.ReadCloser, http.Header, int, error) {

	var urlstr string

	if params.Encode() != "" {
		urlstr = c.restURL.String() + endpoint + "?" + params.Encode()

	} else {
		urlstr = c.restURL.String() + endpoint
	}

	//If Log enabled log URL
	logDebug(ctx, c, fmt.Sprintf("Request Url: %v, Method: %v, PxSession-ID: %v", urlstr, method, PxSessionId))

	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
	case http.MethodDelete, http.MethodGet, http.MethodOptions:
	default:
		return nil, nil, 0, &PxError{Message: fmt.Sprintf("Method is not recognised: %s", method)}
	}

	var body *bytes.Buffer

	// If is File -> no encoding
	if isFile {
		body = bytes.NewBuffer(data.([]byte))

		// PROFFIX REST API Bugfix: Complains if no empty JSON {} is sent
	} else if data == nil || data == "" {

		// If data is emtpy or nil -> send empty JSON Object
		var b bytes.Buffer
		b.Write([]byte("{}"))

		body = &b
		// If not nil -> encode data and send as JSON
	} else {

		body = new(bytes.Buffer)
		encoder := json.NewEncoder(body)

		if err := encoder.Encode(data); err != nil {
			return nil, nil, 0, &PxError{Message: fmt.Sprintf("JSON Encoding failed: %s", err)}
		}
	}
	// End of PROFFIX REST-API Bugfix

	req, err := http.NewRequest(method, urlstr, body)
	if req == nil {
		return nil, nil, 0, err
	}
	req = req.WithContext(ctx)

	if err != nil {
		return nil, nil, 0, err
	}

	// Remove JSON Header is Request is file
	if !isFile {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set PxSessionId in Header
	req.Header.Set("pxsessionid", PxSessionId)

	resp, err := c.client.Do(req)

	if resp != nil && (resp.StatusCode >= 300 || resp.StatusCode < 200) {

		if resp != nil {
			return nil, nil, resp.StatusCode, NewPxError(resp.Body, resp.StatusCode, endpoint)
		} else {

			return nil, nil, 0, &PxError{}
		}

	}

	if resp != nil {
		logDebug(ctx, c, fmt.Sprintf("Response Url: %v, Method: %v, PxSession-ID: %v Status: %v", urlstr, method, PxSessionId, resp.StatusCode))

		// Update the PxSessionId
		c.updatePxSessionId(resp.Header)

		return resp.Body, resp.Header, resp.StatusCode, &PxError{}

		// If everything fails -> logout
	} else {

		_, _ = c.Logout(ctx)
		return nil, nil, 0, &PxError{Message: fmt.Sprintf("%v", err)}
	}

}

// Post Request for PROFFIX REST-API
// Accepts Context, Endpoint and Data as Input
// Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Post(ctx context.Context, endpoint string, data interface{}) (io.ReadCloser, http.Header, int, error) {
	err := c.Login(ctx)
	if !err.(*PxError).isNull() {

		return nil, nil, 0, &PxError{Message: fmt.Sprintf("%v", err)}
	}
	request, header, statuscode, err := c.request(ctx, "POST", endpoint, url.Values{}, false, data)

	//If Log enabled in options log data
	logDebug(ctx, c, fmt.Sprintf("Sent data in POST-Request: %v", data))

	if !err.(*PxError).isNull() {
		return request, header, statuscode, err
	} else {
		return request, header, statuscode, nil
	}

}

// Put Request for PROFFIX REST-API
// Accepts Endpoint and Data as Input
// Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Put(ctx context.Context, endpoint string, data interface{}) (io.ReadCloser, http.Header, int, error) {
	err := c.Login(ctx)
	if !err.(*PxError).isNull() {
		return nil, nil, 0, err
	}
	request, header, statuscode, err := c.request(ctx, "PUT", endpoint, url.Values{}, false, data)

	//If Log enabled in options log data
	logDebug(ctx, c, fmt.Sprintf("Sent data in PUT-Request: %v", data))

	if !err.(*PxError).isNull() {
		return request, header, statuscode, err
	} else {
		return request, header, statuscode, nil
	}
}

// Get Request for PROFFIX REST-API
// Accepts Endpoint and url.Values as Input
// Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Get(ctx context.Context, endpoint string, params url.Values) (io.ReadCloser, http.Header, int, error) {

	err := c.Login(ctx)

	if !err.(*PxError).isNull() {
		return nil, nil, 0, err
	}

	request, header, statuscode, err := c.request(ctx, "GET", endpoint, params, false, nil)
	if !err.(*PxError).isNull() {
		return request, header, statuscode, err
	}

	return request, header, statuscode, nil

}

// Patch Request for PROFFIX REST-API
// Accepts Context, Endpoint and Data as Input
// Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Patch(ctx context.Context, endpoint string, data interface{}) (io.ReadCloser, http.Header, int, error) {
	err := c.Login(ctx)
	if !err.(*PxError).isNull() {
		return nil, nil, 0, err
	}
	request, header, statuscode, err := c.request(ctx, "PATCH", endpoint, url.Values{}, false, data)

	//If Log enabled in options log data
	logDebug(ctx, c, fmt.Sprintf("Sent data in PATCH-Request: %v", data))

	if !err.(*PxError).isNull() {
		return request, header, statuscode, err
	} else {
		return request, header, statuscode, nil
	}
}

// Delete Request for PROFFIX REST-API
// Accepts Endpoint and url.Values as Input
// Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Delete(ctx context.Context, endpoint string) (io.ReadCloser, http.Header, int, error) {
	err := c.Login(ctx)
	if !err.(*PxError).isNull() {
		return nil, nil, 0, err
	}

	request, header, statuscode, err := c.request(ctx, "DELETE", endpoint, nil, false, nil)

	if !err.(*PxError).isNull() {
		return request, header, statuscode, err
	} else {
		return request, header, statuscode, nil
	}
}

// Info Request for PROFFIX REST-API
// Accepts Webservice Key as Input or if left empty uses key from options parameter
// Returns Info about API as io.ReadCloser,error
func (c *Client) Info(ctx context.Context, pxapi string) (io.ReadCloser, error) {

	var endpoint = "PRO/Info"
	param := url.Values{}

	//If no Key submitted in Function use Options
	if pxapi == "" {
		//Set Key from Options
		param.Set("key", c.option.Key)
	} else {
		param.Set("key", pxapi)
	}

	request, _, _, err := c.request(ctx, "GET", endpoint, param, false, "")

	if !err.(*PxError).isNull() {
		return request, err
	} else {
		return request, nil
	}
}

// Database Request for PROFFIX REST-API
// Accepts Webservice Key as Input or if left empty uses key from options parameter
// Returns Database Info as io.ReadCloser,error
func (c *Client) Database(ctx context.Context, pxapi string) (io.ReadCloser, error) {

	param := url.Values{}

	//If no Key submitted in Function use Options
	if pxapi == "" {
		//Set Key from Options
		param.Set("key", c.option.Key)
	} else {
		param.Set("key", pxapi)
	}

	request, _, _, err := c.request(ctx, "GET", "PRO/Datenbank", param, false, "")

	if !err.(*PxError).isNull() {
		return request, err
	} else {
		return request, nil
	}
}

// GetPxSessionId returns the latest PxSessionId from PROFFIX REST-API
func (c *Client) GetPxSessionId() (pxsessionid string) {
	return PxSessionId
}

// logDebug does Log Output if enabled in options
func logDebug(ctx context.Context, c *Client, logtext string) {
	//If Log enabled in options
	if c.option.Log == true {
		log.Print(logtext)
	}
}
