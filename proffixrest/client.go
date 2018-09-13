package proffixrest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	Version = "1.2.0"
)

//Var for storing PxSessionId
var Pxsessionid string

//Client Struct
type Client struct {
	restURL   *url.URL
	Benutzer  string
	Passwort  string
	Datenbank string
	Module    []string
	option    *Options
	rawClient *http.Client
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

//Building new Client
func NewClient(RestURL, apiUser string, apiPassword string, apiDatabase string, apiModule []string, options *Options) (*Client, error) {
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

	path := options.APIPrefix + options.Version + "/"
	restURL.Path = path

	rawClient := &http.Client{}
	rawClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: options.VerifySSL},
	}
	return &Client{
		restURL:   restURL,
		Benutzer:  apiUser,
		Passwort:  apiPassword,
		Datenbank: apiDatabase,
		Module:    apiModule,
		option:    options,
		rawClient: rawClient,
	}, nil
}

//Function for creating new PxSessionId
func (c *Client) createNewPxSessionId() (sessionid string, err error) {

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
	resp, err := c.rawClient.Do(req)
	if err != nil {
		return "", err
	}

	//If Response gives errors print also Body
	if resp.StatusCode == http.StatusBadRequest ||
		resp.StatusCode == http.StatusNotFound ||
		resp.StatusCode == http.StatusInternalServerError {

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		response := buf.String()

		return "", fmt.Errorf("Error on Login: %s", response)
	}
	//If everything OK just return pxsessionid
	return resp.Header.Get("pxsessionid"), nil

}

//Function for Login
//Helper Function for Login to PROFFIX REST-API
//Login is automatically done
func (c *Client) Login() (string, error) {

	// If Pxsessionid doesnt yet exists create a new one
	if Pxsessionid == "" {
		sessionid, err := c.createNewPxSessionId()
		return sessionid, err
		// If Pxsessionid already exists return stored value
	} else {
		return Pxsessionid, nil
	}
}

//LOGOUT Request for PROFFIX REST-API
//Accepts PxSession ID as Input or if left empty uses SessionId from active Session
//Returns Statuscode,error
func (c *Client) Logout() (int, error) {

	_, _, statuscode, err := c.request("DELETE", c.option.LoginEndpoint, url.Values{}, Pxsessionid, nil)
	return statuscode, err
}

//Request Method
//Building the Request Method for Client
func (c *Client) request(method, endpoint string, params url.Values, pxsessionid string, data interface{}) (io.ReadCloser, http.Header, int, error) {

	var urlstr string

	if params.Encode() != "" {
		urlstr = c.restURL.String() + endpoint + "?" + params.Encode()

	} else {
		urlstr = c.restURL.String() + endpoint
	}

	switch method {
	case http.MethodPost, http.MethodPut:
	case http.MethodDelete, http.MethodGet, http.MethodOptions:
	default:
		return nil, nil, 0, fmt.Errorf("Method is not recognised: %s", method)
	}

	body := new(bytes.Buffer)
	encoder := json.NewEncoder(body)
	if err := encoder.Encode(data); err != nil {
		return nil, nil, 0, err
	}

	req, err := http.NewRequest(method, urlstr, body)

	if err != nil {
		return nil, nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("pxsessionid", pxsessionid)
	resp, err := c.rawClient.Do(req)
	if err != nil {
		return nil, nil, resp.StatusCode, err
	}

	return resp.Body, resp.Header, resp.StatusCode, err

}

//POST Request for PROFFIX REST-API
//Accepts Endpoint and Data as Input
//Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Post(endpoint string, data interface{}) (io.ReadCloser, http.Header, int, error) {
	sessionid, err := c.Login()
	if err != nil {
		return nil, nil, 0, err
	}
	request, header, statuscode, err := c.request("POST", endpoint, url.Values{}, sessionid, data)

	//If Login is invalid - try again
	if statuscode == 401 {
		//Get new pxsessionid and write to var
		Pxsessionid, err = c.createNewPxSessionId()
		request, header, statuscode, err := c.request("POST", endpoint, url.Values{}, sessionid, data)

		return request, header, statuscode, err
	}

	//Write the latest pxsessionid back to var
	Pxsessionid = header.Get("pxsessionid")

	return request, header, statuscode, err
}

//PUT Request for PROFFIX REST-API
//Accepts Endpoint and Data as Input
//Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Put(endpoint string, data interface{}) (io.ReadCloser, http.Header, int, error) {
	sessionid, err := c.Login()
	if err != nil {
		return nil, nil, 0, err
	}
	request, header, statuscode, err := c.request("PUT", endpoint, url.Values{}, sessionid, data)

	//If Login is invalid - try again
	if statuscode == 401 {
		//Get new pxsessionid and write to var
		Pxsessionid, err = c.createNewPxSessionId()
		request, header, statuscode, err := c.request("PUT", endpoint, url.Values{}, sessionid, data)

		return request, header, statuscode, err
	}

	//Write the latest pxsessionid back to var
	Pxsessionid = header.Get("pxsessionid")
	//c.Logout(pxsessionid)

	return request, header, statuscode, err
}

//GET Request for PROFFIX REST-API
//Accepts Endpoint and url.Values as Input
//Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Get(endpoint string, params url.Values) (io.ReadCloser, http.Header, int, error) {
	sessionid, err := c.Login()

	if err != nil {
		return nil, nil, 0, err
	}

	request, header, statuscode, err := c.request("GET", endpoint, params, sessionid, nil)

	//If Login is invalid - try again
	if statuscode == 401 {
		//Get new pxsessionid and write to var
		Pxsessionid, err = c.createNewPxSessionId()
		request, header, statuscode, err := c.request("GET", endpoint, params, sessionid, nil)

		return request, header, statuscode, err
	}

	//Write the latest pxsessionid back to var
	Pxsessionid = header.Get("pxsessionid")

	return request, header, statuscode, err
}

//DELETE Request for PROFFIX REST-API
//Accepts Endpoint and url.Values as Input
//Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Delete(endpoint string) (io.ReadCloser, http.Header, int, error) {
	sessionid, err := c.Login()
	if err != nil {
		return nil, nil, 0, err
	}
	request, header, statuscode, err := c.request("DELETE", endpoint, nil, sessionid, nil)

	//If Login is invalid - try again
	if statuscode == 401 {
		//Get new pxsessionid and write to var
		Pxsessionid, err = c.createNewPxSessionId()
		request, header, statuscode, err := c.request("DELETE", endpoint, nil, sessionid, nil)

		return request, header, statuscode, err
	}

	//Write the latest pxsessionid back to var
	Pxsessionid = header.Get("pxsessionid")

	return request, header, statuscode, err
}

//INFO Request for PROFFIX REST-API
//Accepts Webservice Key as Input or if left empty uses key from options parameter
//Returns Info about API as io.ReadCloser,error
func (c *Client) Info(pxapi string) (io.ReadCloser, error) {

	param := url.Values{}

	//If no Key submitted in Function use Options
	if pxapi == "" {
		//Set Key from Options
		param.Set("key", c.option.Key)
	} else {
		param.Set("key", pxapi)
	}

	request, _, _, err := c.request("GET", "PRO/Info", param, "", nil)

	return request, err
}

//DATABASE Request for PROFFIX REST-API
//Accepts Webservice Key as Input or if left empty uses key from options parameter
//Returns Database Info as io.ReadCloser,error
func (c *Client) Database(pxapi string) (io.ReadCloser, error) {

	param := url.Values{}

	//If no Key submitted in Function use Options
	if pxapi == "" {
		//Set Key from Options
		param.Set("key", c.option.Key)
	} else {
		param.Set("key", pxapi)
	}

	request, _, _, err := c.request("GET", "PRO/Datenbank", param, "", nil)

	return request, err
}
