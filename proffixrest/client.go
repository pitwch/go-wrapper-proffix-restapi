package proffixrest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

//Version of Wrapper
const (
	Version = "1.6.8"
)

//Var for storing PxSessionId
//Can be set manually in case of need or for testing
var Pxsessionid string

// DefaultHTTPTransport is an http.RoundTripper that has DisableKeepAlives set true.
var DefaultHTTPTransport = &http.Transport{
	DisableKeepAlives: true,
}

// DefaultHTTPClient is an http.Client with the DefaultHTTPTransport and (Cookie) Jar set nil.
var DefaultHTTPClient = &http.Client{
	Jar:       nil,
	Transport: DefaultHTTPTransport,
}

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
func (c *Client) Login(ctx context.Context) (string, error) {

	// If Pxsessionid doesnt yet exists create a new one
	if Pxsessionid == "" {
		sessionid, err := c.createNewPxSessionId(ctx)

		return sessionid, err
		// If Pxsessionid already exists return stored value
	} else {
		return Pxsessionid, nil
	}
}

//LOGOUT Request for PROFFIX REST-API
//Accepts PxSession ID as Input or if left empty uses SessionId from active Session
//Returns Statuscode,error
func (c *Client) Logout(ctx context.Context) (int, error) {
	//Just logout if we have a valid PxSessionid
	if Pxsessionid != "" {

		//Delete Login Object from PROFFIX REST-API
		_, _, statuscode, err := c.request(ctx, "DELETE", c.option.LoginEndpoint, url.Values{}, Pxsessionid, nil)

		//Set PxSessionId to ""
		Pxsessionid = ""
		return statuscode, err
	} else {
		return 0, nil
	}

}

//Request Method
//Building the Request Method for Client
func (c *Client) request(ctx context.Context, method, endpoint string, params url.Values, pxsessionid string, data interface{}) (io.ReadCloser, http.Header, int, error) {

	var urlstr string

	if params.Encode() != "" {
		urlstr = c.restURL.String() + endpoint + "?" + params.Encode()

	} else {
		urlstr = c.restURL.String() + endpoint
	}

	//If Log enabled log URL
	logDebug(ctx, c, fmt.Sprintf("Request-URL: %v, Method: %v, PxSession-ID: %v", urlstr, method, pxsessionid))

	switch method {
	case http.MethodPost, http.MethodPut:
	case http.MethodDelete, http.MethodGet, http.MethodOptions:
	default:
		return nil, nil, 0, fmt.Errorf("Method is not recognised: %s", method)
	}

	body := new(bytes.Buffer)
	encoder := json.NewEncoder(body)
	if err := encoder.Encode(data); err != nil {
		return nil, nil, 0, fmt.Errorf("JSON Encoding failed: %s", err)
	}

	req, err := http.NewRequest(method, urlstr, body)
	req = req.WithContext(ctx)

	if err != nil {
		return nil, nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("pxsessionid", pxsessionid)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, resp.StatusCode, err
	}

	return resp.Body, resp.Header, resp.StatusCode, err

}

//POST Request for PROFFIX REST-API
//Accepts Context, Endpoint and Data as Input
//Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Post(ctx context.Context, endpoint string, data interface{}) (io.ReadCloser, http.Header, int, error) {
	sessionid, err := c.Login(ctx)
	if err != nil {
		return nil, nil, 0, err
	}
	request, header, statuscode, err := c.request(ctx, "POST", endpoint, url.Values{}, sessionid, data)

	//If Log enabled in options log data
	logDebug(ctx, c, fmt.Sprintf("Sent data in POST-Request: %v", data))

	//If Login is invalid - try again
	if statuscode == 401 {
		//Get new pxsessionid and write to var
		Pxsessionid, err = c.createNewPxSessionId(ctx)
		request, header, statuscode, err := c.request(ctx, "POST", endpoint, url.Values{}, sessionid, data)

		return request, header, statuscode, err
	}

	//If Statuscode not 201
	if statuscode != 201 {
		return request, header, statuscode, errorFormatterPx(ctx, c, statuscode, request)
	}

	//Write the latest pxsessionid back to var
	Pxsessionid = header.Get("pxsessionid")

	return request, header, statuscode, err
}

//PUT Request for PROFFIX REST-API
//Accepts Endpoint and Data as Input
//Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Put(ctx context.Context, endpoint string, data interface{}) (io.ReadCloser, http.Header, int, error) {
	sessionid, err := c.Login(ctx)
	if err != nil {
		return nil, nil, 0, err
	}
	request, header, statuscode, err := c.request(ctx, "PUT", endpoint, url.Values{}, sessionid, data)

	//If Log enabled in options log data
	logDebug(ctx, c, fmt.Sprintf("Sent data in PUT-Request: %v", data))

	//If Login is invalid - try again
	if statuscode == 401 {
		//Get new pxsessionid and write to var
		Pxsessionid, err = c.createNewPxSessionId(ctx)
		request, header, statuscode, err := c.request(ctx, "PUT", endpoint, url.Values{}, sessionid, data)

		return request, header, statuscode, err
	}

	//If Statuscode not 204
	if statuscode != 204 {
		return request, header, statuscode, errorFormatterPx(ctx, c, statuscode, request)
	}

	//Write the latest pxsessionid back to var
	Pxsessionid = header.Get("pxsessionid")
	//c.Logout(pxsessionid)

	return request, header, statuscode, err
}

//GET Request for PROFFIX REST-API
//Accepts Endpoint and url.Values as Input
//Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Get(ctx context.Context, endpoint string, params url.Values) (io.ReadCloser, http.Header, int, error) {
	sessionid, err := c.Login(ctx)

	if err != nil {
		return nil, nil, 0, err
	}

	request, header, statuscode, err := c.request(ctx, "GET", endpoint, params, sessionid, nil)

	//If Login is invalid - try again
	if statuscode == 401 {
		//Get new pxsessionid and write to var
		Pxsessionid, err = c.createNewPxSessionId(ctx)
		request, header, statuscode, err := c.request(ctx, "GET", endpoint, params, sessionid, nil)

		return request, header, statuscode, err
	}

	//If Statuscode not 200
	if statuscode != 200 {
		return request, header, statuscode, errorFormatterPx(ctx, c, statuscode, request)
	}

	//Write the latest pxsessionid back to var
	Pxsessionid = header.Get("pxsessionid")

	return request, header, statuscode, err
}

//DELETE Request for PROFFIX REST-API
//Accepts Endpoint and url.Values as Input
//Returns io.ReadCloser,http.Header,Statuscode,error
func (c *Client) Delete(ctx context.Context, endpoint string) (io.ReadCloser, http.Header, int, error) {
	sessionid, err := c.Login(ctx)
	if err != nil {
		return nil, nil, 0, err
	}
	request, header, statuscode, err := c.request(ctx, "DELETE", endpoint, nil, sessionid, nil)

	//If Login is invalid - try again
	if statuscode == 401 {
		//Get new pxsessionid and write to var
		Pxsessionid, err = c.createNewPxSessionId(ctx)
		request, header, statuscode, err := c.request(ctx, "DELETE", endpoint, nil, sessionid, nil)

		return request, header, statuscode, err
	}

	//If Statuscode not 204
	if statuscode != 204 {
		return request, header, statuscode, errorFormatterPx(ctx, c, statuscode, request)
	}

	//Write the latest pxsessionid back to var
	Pxsessionid = header.Get("pxsessionid")

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

	request, _, _, err := c.request(ctx, "GET", "PRO/Info", param, "", nil)

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

	request, _, _, err := c.request(ctx, "GET", "PRO/Datenbank", param, "", nil)

	return request, err
}

func GetFiltererCount(header http.Header) (total int) {
	type PxMetadata struct {
		FilteredCount int
	}
	var pxmetadata PxMetadata
	head := header.Get("pxmetadata")
	json.Unmarshal([]byte(head), &pxmetadata)

	return pxmetadata.FilteredCount
}

func (c *Client) GetBatch(ctx context.Context, endpoint string, params url.Values, batchsize int) (result []byte, total int, err error) {

	//Create collector for collecting requests
	var collector []byte

	//If batchsize not set in function / is 0
	if batchsize == 0 {
		//Set it to default
		batchsize = c.option.Batchsize
	}

	//If params were set to nil set to empty url.Values{}
	if params == nil {
		params = url.Values{}
	}

	//Create Setting Query for setting up the sync and getting the FilteredCount Header

	//Store original params in paramquery
	paramquery := params

	//Delete Limit from params in case user defined it
	paramquery.Del("Limit")

	//Add Limit / Batchsize to params
	paramquery.Add("Limit", cast.ToString(batchsize))

	//Query Endpoint for results
	rc, header, status, err := c.Get(ctx, endpoint, params)

	if err != nil {
		return nil, 0, err
	}

	//Read from rc into settingResp
	settingResp, err := ioutil.ReadAll(rc)

	//Put setting query into collector
	collector = append(collector, settingResp[:]...)

	//If Status not 200 log Error message
	if status != 200 {
		return nil, 0, err
	}
	//Get total available objects
	totalEntries := GetFiltererCount(header)

	//Unmarshall to interface so we can count elements in setting query
	var jsonObjs interface{}
	err = json.Unmarshal([]byte(collector), &jsonObjs)

	firstQueryCount := len(jsonObjs.([]interface{}))

	//Add firstQueryCount to the totalEntriesCount so we can save queries
	totalEntriesCount := firstQueryCount

	//If total available objects are already all requested --> return and exit
	if totalEntriesCount < totalEntries {

		//Loop over requests until we have all available objects
		for totalEntriesCount < totalEntries {

			//Reset the params to original
			paramquery := url.Values{}
			paramquery = params

			//To be sure; remove old limit + offset
			paramquery.Del("Limit")
			paramquery.Del("Offset")

			//Add new limit + offset
			paramquery.Add("Limit", cast.ToString(batchsize))
			paramquery.Add("Offset", cast.ToString(totalEntriesCount))

			//Fire request with new limit + offset
			bc, _, _, _ := c.Get(ctx, endpoint, paramquery)

			//Read from bc into temporary batchResp
			batchResp, err := ioutil.ReadAll(bc)

			collector = append(collector, batchResp[:]...)
			if err != nil {
				panic(err)
			}

			//If Status not 200 log Error message
			if status != 200 {
				panic(string(collector))
			}

			//Unmarshall to interface so we can count elements in totalEntriesCount query
			var jsonObjs2 interface{}
			json.Unmarshal([]byte(batchResp), &jsonObjs2)
			totalEntriesCount += len(jsonObjs2.([]interface{}))
		}

		//Clean the collector, replace JSON Tag to get one big JSON array
		cleanCollector := bytes.Replace(collector, []byte("]["), []byte(","), -1)

		return cleanCollector, totalEntriesCount, err

		//If count of elements in first query is bigger or same than FilteredCount return
	} else {
		//Return
		return collector, totalEntriesCount, err
	}
	//Return
	return collector, totalEntriesCount, err
}

func errorFormatterPx(ctx context.Context, c *Client, statuscode int, request io.Reader) (err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(request)
	errbyte := buf.Bytes()
	errstr := buf.String()

	//Do Logout...
	logoutstatus, err := c.Logout(ctx)
	logDebug(ctx, c, fmt.Sprintf("Logout after Error with Status - Code: %v", logoutstatus))

	//Define Error Struct
	parsedError := ErrorStruct{}

	//Try to parse JSON in ErrorStruct
	err = json.Unmarshal(errbyte, &parsedError)

	//If error on parse return plain text
	if err != nil {
		return fmt.Errorf("ERROR: %v", errstr)
	}
	return fmt.Errorf("Error Statuscode %v Type %s: %s", statuscode, parsedError.Type, parsedError.Message)
}

func logDebug(ctx context.Context, c *Client, logtext string) {
	//If Log enabled in options
	if c.option.Log == true {
		log.Print(logtext)
	}
}
