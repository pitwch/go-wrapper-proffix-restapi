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
	Version = "1.0.0"
)

type Client struct {
	restURL   *url.URL
	Benutzer  string
	Passwort  string
	Datenbank string
	Module    []string
	option    *Options
	rawClient *http.Client
}

type LoginStruct struct {
	Benutzer  string         `json:"Benutzer"`
	Passwort  string         `json:"Passwort"`
	Datenbank DatabaseStruct `json:"Datenbank"`
	Module    []string       `json:"Module"`
}

type DatabaseStruct struct {
	Name string `json:"Name"`
}

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
func (c *Client) PxLogin() (string, error) {
	urlstr := c.restURL.String() + c.option.LoginEndpoint

	data := LoginStruct{c.Benutzer, c.Passwort, DatabaseStruct{c.Datenbank}, c.Module}
	body := new(bytes.Buffer)
	encoder := json.NewEncoder(body)
	if err := encoder.Encode(data); err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", urlstr, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.rawClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusBadRequest ||
		resp.StatusCode == http.StatusUnauthorized ||
		resp.StatusCode == http.StatusNotFound ||
		resp.StatusCode == http.StatusInternalServerError {

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		newStr := buf.String()

		fmt.Printf(newStr)

		return "", fmt.Errorf("Request failed: %s", resp.Status)
	}
	return resp.Header.Get("PxSessionId"), nil
}

func (c *Client) PxLogout(pxsessionid string) (int, error) {
	urlstr := c.restURL.String() + c.option.LoginEndpoint

	req, err := http.NewRequest("DELETE", urlstr, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PxSessionId", pxsessionid)
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
func (c *Client) request(method, endpoint string, params url.Values, pxsessionid string, data interface{}) (io.ReadCloser, http.Header, error) {

	var urlstr string

	if params.Encode() != "" {
		urlstr = c.restURL.String() + endpoint + "?" + params.Encode()

	} else {
		urlstr = c.restURL.String() + endpoint
	}

	log.Print(urlstr)
	//if c.restURL.Scheme == "https" {
	//	urlstr += "?" + c.basicAuth(params)
	//} else {
	//	urlstr += "?" + c.oauth(method, urlstr, params)
	//}
	switch method {
	case http.MethodPost, http.MethodPut:
	case http.MethodDelete, http.MethodGet, http.MethodOptions:
	default:
		return nil, nil, fmt.Errorf("Method is not recognised: %s", method)
	}

	body := new(bytes.Buffer)
	encoder := json.NewEncoder(body)
	if err := encoder.Encode(data); err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(method, urlstr, body)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PxSessionId", pxsessionid)
	resp, err := c.rawClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode == http.StatusBadRequest ||
		resp.StatusCode == http.StatusUnauthorized ||
		resp.StatusCode == http.StatusNotFound ||
		resp.StatusCode == http.StatusInternalServerError {

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		body := buf.String()
		log.Print(body)

		return nil, nil, fmt.Errorf("Request failed: %s", resp.Status)
	}
	return resp.Body, resp.Header, nil
}

func (c *Client) Post(endpoint string, data interface{}) (io.ReadCloser, http.Header, error) {
	pxsessionid, err := c.PxLogin()
	if err != nil {
		return nil, nil, err
	}
	request, header, err := c.request("POST", endpoint, url.Values{}, pxsessionid, data)

	c.PxLogout(pxsessionid)

	return request, header, err
}

func (c *Client) Put(endpoint string, data interface{}) (io.ReadCloser, http.Header, error) {
	pxsessionid, err := c.PxLogin()
	if err != nil {
		return nil, nil, err
	}
	request, header, err := c.request("PUT", endpoint, url.Values{}, pxsessionid, data)

	c.PxLogout(pxsessionid)

	return request, header, err
}

func (c *Client) Get(endpoint string, params url.Values) (io.ReadCloser, http.Header, error) {
	pxsessionid, err := c.PxLogin()
	if err != nil {
		return nil, nil, err
	}
	request, header, err := c.request("GET", endpoint, params, pxsessionid, nil)

	c.PxLogout(pxsessionid)

	return request, header, err
}

func (c *Client) Delete(endpoint string, params url.Values) (io.ReadCloser, http.Header, error) {
	pxsessionid, err := c.PxLogin()
	if err != nil {
		return nil, nil, err
	}
	request, header, err := c.request("DELETE", endpoint, params, pxsessionid, nil)

	c.PxLogout(pxsessionid)

	return request, header, err
}

func (c *Client) Info(pxapi string) (io.ReadCloser, error) {

	param := url.Values{}
	param.Set("key", pxapi)
	request, _, err := c.request("GET", "PRO/Info", param, "", nil)

	return request, err
}

func (c *Client) Database(pxapi string) (io.ReadCloser, error) {

	param := url.Values{}
	param.Set("key", pxapi)
	request, _, err := c.request("GET", "PRO/Datenbank", param, "", nil)

	return request, err
}
