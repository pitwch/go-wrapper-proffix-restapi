package proffixrest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"log"
)

const (
	Version       = "1.0.0"
	HashAlgorithm = "HMAC-SHA256"
)

type Client struct {
	restURL   *url.URL
	Benutzer  string
	Passwort  string
	Datenbank string
	Module    string
	option    *Options
	rawClient *http.Client
}

type LoginStruct struct {
	Benutzer  string         `json:"Benutzer"`
	Passwort  string         `json:"Passwort"`
	Datenbank DatabaseStruct `json:"Datenbank"`
	Module    string         `json:"Module"`
}

type DatabaseStruct struct {
	Name        string `json:"Name"`
	Bezeichnung string `json:"Bezeichnung"`
}

func NewClient(restapi, Benutzer string, Passwort string, Datenbank string, Module string, option *Options) (*Client, error) {
	restURL, err := url.Parse(restapi)
	if err != nil {
		return nil, err
	}

	if option == nil {
		option = &Options{}
	}

	if option.Version == "" {
		option.Version = "v2"
	}

	if option.APIPrefix == "" {
		option.APIPrefix = "/pxapi/"
	}

	if option.LoginEndpoint == "" {
		option.LoginEndpoint = "PRO/Login"
	}

	if option.UserAgent == "" {
		option.UserAgent = "go-wrapper-proffix-restapi " + Version
	}

	path := option.APIPrefix + option.Version + "/"
	restURL.Path = path

	rawClient := &http.Client{}
	rawClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: option.VerifySSL},
	}
	return &Client{
		restURL:   restURL,
		Benutzer:  Benutzer,
		Passwort:  Passwort,
		Datenbank: Datenbank,
		Module:    Module,
		option:    option,
		rawClient: rawClient,
	}, nil
}

func (c *Client) PxLogin(option *Options) string {

	loginJson := LoginStruct{c.Benutzer, c.Passwort, DatabaseStruct{c.Datenbank,""}, c.Module}
	jsonStr, err := json.Marshal(loginJson)
	if err != nil {
		panic(err)
	}
	log.Println(loginJson)
	req, err := http.NewRequest("POST", "https://fernwartung.pitw.ch:955/pxapi/v2/PRO/Login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		log.Println("Error", resp.Body)

	}

	return string(resp.Header.Get("PxSessionId"))
}

func (c *Client) request(method, endpoint string, params url.Values, data interface{}) (io.ReadCloser, error) {
	urlstr := c.restURL.String() + endpoint
	if params == nil {
		params = make(url.Values)
	}
	//if c.restURL.Scheme == "https" {
	//	urlstr += "?" + c.basicAuth(params)
	//} else {
	//	urlstr += "?" + c.oauth(method, urlstr, params)
	//}
	switch method {
	case http.MethodPost, http.MethodPut:
	case http.MethodDelete, http.MethodGet, http.MethodOptions:
	default:
		return nil, fmt.Errorf("Method is not recognised: %s", method)
	}

	body := new(bytes.Buffer)
	encoder := json.NewEncoder(body)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	log.Println(urlstr)

	req, err := http.NewRequest(method, urlstr, body)
	if err != nil {
		return nil, err
	}
	x := c.PxLogin(c.option)
	log.Println(x)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PxSessionId", x)
	resp, err := c.rawClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusBadRequest ||
		resp.StatusCode == http.StatusUnauthorized ||
		resp.StatusCode == http.StatusNotFound ||
		resp.StatusCode == http.StatusInternalServerError {

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		newStr := buf.String()

		fmt.Printf(newStr)

		return nil, fmt.Errorf("Request failed: %s", resp.Status)
	}
	return resp.Body, nil
}

func (c *Client) Post(endpoint string, data interface{}) (io.ReadCloser, error) {
	return c.request("POST", endpoint, nil, data)
}

func (c *Client) Put(endpoint string, data interface{}) (io.ReadCloser, error) {
	return c.request("PUT", endpoint, nil, data)
}

func (c *Client) Get(endpoint string, params url.Values) (io.ReadCloser, error) {
	return c.request("GET", endpoint, params, nil)
}

func (c *Client) Delete(endpoint string, params url.Values) (io.ReadCloser, error) {
	return c.request("POST", endpoint, params, nil)
}

func (c *Client) Options(endpoint string) (io.ReadCloser, error) {
	return c.request("OPTIONS", endpoint, nil, nil)
}
