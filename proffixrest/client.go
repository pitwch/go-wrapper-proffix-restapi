package proffixrest

import (
	"bytes"
	"context"
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
	Module    []string       `json:"Module"`
}

//Database Struct
type DatabaseStruct struct {
	Name string `json:"Name"`
}

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

//Building new Client
func NewClient(RestURL string, apiUser string, apiPassword string, apiDatabase string, apiModule []string, options *Options) (*Client, error) {
	restURL, err := url.Parse(RestURL)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &Options{}
	}

	if options.Version == "" {
		options.Version = "v2"
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

	//Set default batchsize for batch requests
	if options.Batchsize == 0 {
		options.Batchsize = 200
	}

	if DefaultHTTPClient == nil {

		DefaultHTTPClient = http.DefaultClient
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

//Function for creating new PxSessionId
func (c *Client) createNewPxSessionId(ctx context.Context) (sessionid string, err error) {

	//Check URL, else exit
	_, err = url.ParseRequestURI(c.restURL.String())
	if err != nil {
		return "", fmt.Errorf("URL in wrong format: %s", err)
	}

	//Build Login URL
	urlstr := c.restURL.String() + c.option.LoginEndpoint

	//Create Login Body
	data := LoginStruct{c.Benutzer, c.Passwort, DatabaseStruct{c.Datenbank}, c.Module}

	//Encode Login Body to JSON
	body := new(bytes.Buffer)
	encoder := json.NewEncoder(body)
	if err := encoder.Encode(data); err != nil {
		return "", err
	}

	//Build Login Request
	req, err := http.NewRequest("POST", urlstr, body)
	if err != nil {
		return "", err
	}

	//Set Login Header
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}

	//If Response gives errors print also Body
	if resp.StatusCode == http.StatusBadRequest ||
		resp.StatusCode == http.StatusNotFound ||
		resp.StatusCode == http.StatusMethodNotAllowed ||
		resp.StatusCode == http.StatusInternalServerError {

		return "", errorFormatterPx(ctx, c, resp.StatusCode, resp.Body)
	}
	//If everything OK just return pxsessionid
	return resp.Header.Get("pxsessionid"), nil

}

//Function for Login
//Helper Function for Login to PROFFIX REST-API
//Login is automatically done
func (c *Client) Login(ctx context.Context) error {

	// If Pxsessionid doesnt yet exists create a new one
	if PxSessionId == "" {
		sessionid, err := c.createNewPxSessionId(ctx)
		PxSessionId = sessionid
		return err
		// If Pxsessionid already exists return stored value
	} else {
		return nil
	}
}

// updatePxSessionId updates the stored PxSessionId
func (c *Client) updatePxSessionId(header http.Header) {

	// Just update if PxSessionId in Header is not empty
	if header.Get("pxsessionid") != "" {
		PxSessionId = header.Get("pxsessionid")
	}

}

//LOGOUT Request for PROFFIX REST-API
//Accepts PxSession ID as Input or if left empty uses SessionId from active Session
//Returns Statuscode,error
func (c *Client) Logout(ctx context.Context) (int, error) {
	//Just logout if we have a valid PxSessionid
	if PxSessionId != "" {

		//Delete Login Object from PROFFIX REST-API
		_, _, statuscode, err := c.request(ctx, "DELETE", c.option.LoginEndpoint, url.Values{}, nil)

		//Set Pxsessionid to empty string
		PxSessionId = ""

		return statuscode, err
	} else {
		return 0, nil
	}

}

//Request Method
//Building the Request Method for Client
func (c *Client) request(ctx context.Context, method, endpoint string, params url.Values, data interface{}) (io.ReadCloser, http.Header, int, error) {

	var urlstr string

	if params.Encode() != "" {
		urlstr = c.restURL.String() + endpoint + "?" + params.Encode()

	} else {
		urlstr = c.restURL.String() + endpoint
	}

	//If Log enabled log URL
	logDebug(ctx, c, fmt.Sprintf("Request Url: %v, Method: %v, PxSession-ID: %v", urlstr, method, PxSessionId))

	switch method {
	case http.MethodPost, http.MethodPut:
	case http.MethodDelete, http.MethodGet, http.MethodOptions:
	default:
		return nil, nil, 0, fmt.Errorf("Method is not recognised: %s", method)
	}

	var body *bytes.Buffer

	//PROFFIX REST API Bugfix: Complains if no empty JSON {} is sent

	//If data is emtpy or nil -> send empty JSON Object
	if data == nil || data == "" {
		var b bytes.Buffer
		b.Write([]byte("{}"))

		body = &b
		//If not nil -> encode data and send as JSON
	} else {

		body = new(bytes.Buffer)
		encoder := json.NewEncoder(body)

		if err := encoder.Encode(data); err != nil {
			return nil, nil, 0, fmt.Errorf("JSON Encoding failed: %s", err)
		}
	}
	//End of PROFFIX REST-API Bugfix

	req, err := http.NewRequest(method, urlstr, body)
	req = req.WithContext(ctx)

	if err != nil {
		return nil, nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	//Set PxSessionId in Header
	req.Header.Set("pxsessionid", PxSessionId)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, resp.StatusCode, err
	}
	logDebug(ctx, c, fmt.Sprintf("Response Url: %v, Method: %v, PxSession-ID: %v Status: %v", urlstr, method, PxSessionId, resp.StatusCode))

	//Update the PxSessionId
	c.updatePxSessionId(resp.Header)

	return resp.Body, resp.Header, resp.StatusCode, err

}

//POST Request for PROFFIX REST-API
//Accepts Context, Endpoint and Data as Input
//Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Post(ctx context.Context, endpoint string, data interface{}) (io.ReadCloser, http.Header, int, error) {
	err := c.Login(ctx)
	if err != nil {
		return nil, nil, 0, err
	}
	request, header, statuscode, err := c.request(ctx, "POST", endpoint, url.Values{}, data)

	//If Log enabled in options log data
	logDebug(ctx, c, fmt.Sprintf("Sent data in POST-Request: %v", data))

	//If Login is invalid - try again
	if statuscode == 401 {
		//Set PxSessionId to empty string
		PxSessionId = ""

		//Get new pxsessionid and write to var
		err = c.Login(ctx)

		//Repeat Request with new SessionId
		request, header, statuscode, err := c.request(ctx, "POST", endpoint, url.Values{}, data)

		return request, header, statuscode, err
	}

	//If Statuscode not 201
	if statuscode != 201 {
		return request, header, statuscode, errorFormatterPx(ctx, c, statuscode, request)
	}

	return request, header, statuscode, err
}

//PUT Request for PROFFIX REST-API
//Accepts Endpoint and Data as Input
//Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Put(ctx context.Context, endpoint string, data interface{}) (io.ReadCloser, http.Header, int, error) {
	err := c.Login(ctx)
	if err != nil {
		return nil, nil, 0, err
	}
	request, header, statuscode, err := c.request(ctx, "PUT", endpoint, url.Values{}, data)

	//If Log enabled in options log data
	logDebug(ctx, c, fmt.Sprintf("Sent data in PUT-Request: %v", data))

	//If Login is invalid - try again
	if statuscode == 401 {
		//Set PxSessionId to empty string
		PxSessionId = ""

		//Get new pxsessionid and write to var
		err = c.Login(ctx)

		//Repeat Request with new SessionId
		request, header, statuscode, err := c.request(ctx, "PUT", endpoint, url.Values{}, data)

		return request, header, statuscode, err
	}

	//If Statuscode not 204
	if statuscode != 204 {
		return request, header, statuscode, errorFormatterPx(ctx, c, statuscode, request)
	}

	return request, header, statuscode, err
}

//GET Request for PROFFIX REST-API
//Accepts Endpoint and url.Values as Input
//Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Get(ctx context.Context, endpoint string, params url.Values) (io.ReadCloser, http.Header, int, error) {

	err := c.Login(ctx)

	if err != nil {
		return nil, nil, 0, err
	}

	request, header, statuscode, err := c.request(ctx, "GET", endpoint, params, nil)

	//If Login is invalid - try again
	if statuscode == 401 {

		//Set PxSessionId to empty string
		PxSessionId = ""

		//Repeat Request with new SessionId
		request, header, statuscode, err := c.request(ctx, "GET", endpoint, url.Values{}, nil)

		return request, header, statuscode, err
	}

	//If Statuscode not 200
	if statuscode != 200 {
		return request, header, statuscode, errorFormatterPx(ctx, c, statuscode, request)
	}

	return request, header, statuscode, err
}

//DELETE Request for PROFFIX REST-API
//Accepts Endpoint and url.Values as Input
//Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Delete(ctx context.Context, endpoint string) (io.ReadCloser, http.Header, int, error) {
	err := c.Login(ctx)
	if err != nil {
		return nil, nil, 0, err
	}
	request, header, statuscode, err := c.request(ctx, "DELETE", endpoint, nil, nil)

	//If Login is invalid - try again
	if statuscode == 401 {

		//Set PxSessionId to empty string
		PxSessionId = ""

		//Get new pxsessionid and write to var
		err = c.Login(ctx)

		//Repeat Request with new SessionId
		request, header, statuscode, err := c.request(ctx, "DELETE", endpoint, url.Values{}, nil)

		return request, header, statuscode, err
	}

	//If Statuscode not 204
	if statuscode != 204 {
		return request, header, statuscode, errorFormatterPx(ctx, c, statuscode, request)
	}

	return request, header, statuscode, err
}

//INFO Request for PROFFIX REST-API
//Accepts Webservice Key as Input or if left empty uses key from options parameter
//Returns Info about API as io.ReadCloser,error
func (c *Client) Info(ctx context.Context, pxapi string) (io.ReadCloser, error) {

	param := url.Values{}

	//If no Key submitted in Function use Options
	if pxapi == "" {
		//Set Key from Options
		param.Set("key", c.option.Key)
	} else {
		param.Set("key", pxapi)
	}

	request, _, _, err := c.request(ctx, "GET", "PRO/Info", param, "")

	return request, err
}

//DATABASE Request for PROFFIX REST-API
//Accepts Webservice Key as Input or if left empty uses key from options parameter
//Returns Database Info as io.ReadCloser,error
func (c *Client) Database(ctx context.Context, pxapi string) (io.ReadCloser, error) {

	param := url.Values{}

	//If no Key submitted in Function use Options
	if pxapi == "" {
		//Set Key from Options
		param.Set("key", c.option.Key)
	} else {
		param.Set("key", pxapi)
	}

	request, _, _, err := c.request(ctx, "GET", "PRO/Datenbank", param, "")

	return request, err
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
