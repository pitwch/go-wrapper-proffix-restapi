package proffixrest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	Version = "1.1.0"
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

	//If Login Response gives errors print also Body
	if resp.StatusCode == http.StatusBadRequest ||
		resp.StatusCode == http.StatusUnauthorized ||
		resp.StatusCode == http.StatusNotFound ||
		resp.StatusCode == http.StatusInternalServerError {

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		response := buf.String()

		fmt.Printf(response)

		return "", fmt.Errorf("Request failed: %s", resp.Status)
	}
	//If everything OK just return pxsessionid
	return resp.Header.Get("pxsessionid"), nil

}

//Function for Login
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

//Function for Logout
func (c *Client) Logout(pxsessionid string) (int, error) {

	urlstr := c.restURL.String() + c.option.LoginEndpoint

	req, err := http.NewRequest("DELETE", urlstr, nil)
	if err != nil {
		return 0, err
	}

	if pxsessionid != "" {
		req.Header.Set("pxsessionid", pxsessionid)
	} else {
		req.Header.Set("PxSessionId", Pxsessionid)
	}

	resp, err := c.rawClient.Do(req)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode == http.StatusBadRequest ||
		resp.StatusCode == http.StatusUnauthorized ||
		resp.StatusCode == http.StatusNotFound ||
		resp.StatusCode == http.StatusInternalServerError {

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		newStr := buf.String()

		fmt.Printf(newStr)

		return 0, fmt.Errorf("Request failed: %s", resp.Status)
	}
	return resp.StatusCode, nil
}

//Request Method
func (c *Client) request(method, endpoint string, params url.Values, pxsessionid string, data interface{}) (io.ReadCloser, http.Header, int, error) {

	var urlstr string

	if params.Encode() != "" {
		urlstr = c.restURL.String() + endpoint + "?" + params.Encode()

	} else {
		urlstr = c.restURL.String() + endpoint
	}

	log.Print(urlstr)

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

//Post
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

	//c.Logout(sessionid)

	return request, header, statuscode, err
}

//Put
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
		request, header, statuscode, err := c.request("GET", endpoint, url.Values{}, sessionid, nil)

		return request, header, statuscode, err
	}

	//Write the latest pxsessionid back to var
	Pxsessionid = header.Get("pxsessionid")
	//c.Logout(pxsessionid)

	return request, header, statuscode, err
}

//Delete
func (c *Client) Delete(endpoint string, params url.Values) (io.ReadCloser, http.Header, int, error) {
	sessionid, err := c.Login()
	if err != nil {
		return nil, nil, 0, err
	}
	request, header, statuscode, err := c.request("DELETE", endpoint, params, sessionid, nil)

	//If Login is invalid - try again
	if statuscode == 401 {
		//Get new pxsessionid and write to var

		Pxsessionid, err = c.createNewPxSessionId()
		request, header, statuscode, err := c.request("DELETE", endpoint, url.Values{}, sessionid, nil)

		return request, header, statuscode, err
	}

	//Write the latest pxsessionid back to var
	Pxsessionid = header.Get("pxsessionid")
	//c.Logout(pxsessionid)

	return request, header, statuscode, err
}

//Endpoint Info
func (c *Client) Info(pxapi string) (io.ReadCloser, error) {

	param := url.Values{}
	param.Set("key", pxapi)
	request, _, _, err := c.request("GET", "PRO/Info", param, "", nil)

	return request, err
}

//Endpoint Database
func (c *Client) Database(pxapi string) (io.ReadCloser, error) {

	param := url.Values{}
	param.Set("key", pxapi)
	request, _, _, err := c.request("GET", "PRO/Datenbank", param, "", nil)

	return request, err
}
