// Package main provides a command-line interface for the PROFFIX REST-API wrapper.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/pitwch/go-wrapper-proffix-restapi/proffixrest"
	"github.com/xiaost/jsonport"
)

var SRVERSION string

type configcmd struct {
	URL      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Module   string `json:"module"`
	APIKey   string `json:"apikey"`
}

func main() {
	// General Config Options (Params)
	resturl := flag.String("url", os.Getenv("PXR_URL"), "URL of PROFFIX REST-API in Format <https://>URL:Port")
	user := flag.String("user", os.Getenv("PXR_USER"), "User of PROFFIX REST-API")
	password := flag.String("password", os.Getenv("PXR_PASSWORD"), "Password of User. Either as SHA256 - Hash or plain text")
	database := flag.String("database", os.Getenv("PXR_DATABASE"), "Database in PROFFIX")
	module := flag.String("module", os.Getenv("PXR_MODULE"), "Module which are needed for Licence")
	apikey := flag.String("apikey", os.Getenv("PXR_APIKEY"), "API Password / Key (just needed for Endpoints Info / Datenbank")

	// General Config Options (Configfile)
	config := flag.String("config", "", "Use Configfile instead Parameters")

	// Specific Config Options
	method := flag.String("method", "", "Method for query API")
	endpoint := flag.String("endpoint", "", "Endpoint from PROFFIX-REST API")
	data := flag.String("data", "", "Data as JSON")

	// Parameters
	limit := flag.String("limit", "1", "Limit")
	format := flag.String("format", "json", "Format")
	field := flag.String("field", "", "Field")
	depth := flag.Int("depth", 0, "Depth")

	showVersion := flag.Bool("version", false, "Shows the version of go-proffix-restapi-wrapper")
	log := flag.Bool("log", false, "Enable Log")
	// updateFile := flag.String("update", "", "updates the version of a certain file")

	flag.Parse()

	if *showVersion {
		fmt.Printf("go-wrapper-proffix-restapi v%s", SRVERSION)
		return
	}

	if *config != "" {
		// Open jsonFile
		jsonFile, err := os.Open(*config)
		// if  os.Open returns an error then handle it
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Successfully loaded config from file")
		}
		// defer the closing of our jsonFile so that we can parse it later on
		defer func() { _ = jsonFile.Close() }()

		var cfg configcmd

		byteValue, _ := ioutil.ReadAll(jsonFile) // Ignore error, handled below

		// Unmarshmal Configfile to config
		_ = json.Unmarshal(byteValue, &cfg) // Ignore error, use zero values on failure

		// Set Vars from file
		resturl = &cfg.URL
		user = &cfg.User
		password = &cfg.Password
		database = &cfg.Database
		module = &cfg.Module
		apikey = &cfg.APIKey
	}

	param := url.Values{}
	param.Add("limit", *limit)

	ctx := context.Background()

	pxrestcmd, err := proffixrest.NewClient(
		*resturl,
		*user,
		*password,
		*database,
		[]string{*module},
		&proffixrest.Options{Key: *apikey, Log: *log},
	)
	if err != nil {
		fmt.Printf("Fehler: %v", err)
	}

	switch strings.ToLower(*method) {
	case "get":
		fmt.Print(get(ctx, pxrestcmd, *endpoint, param, *format, *depth, *field))
		// Logout
		_, _ = pxrestcmd.Logout(ctx)
	case "post":
		fmt.Print(post(ctx, pxrestcmd, *endpoint, *data))
		// Logout
		_, _ = pxrestcmd.Logout(ctx)
	}

}
func responseFormatter(closer io.ReadCloser, format string, depth int, field string) interface{} {
	// Buffer decode for plain text response
	if closer != nil {

		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(closer)

		respStr := buf.String()
		respByte := buf.Bytes()
		defer func() { _ = closer.Close() }()

		switch strings.ToLower(format) {
		case "json":
			return respStr
		case "string":
			j, _ := jsonport.Unmarshal(respByte, depth, field) // Ignore error, handle nil below
			switch {
			case j.IsString():
				str, err := j.GetString()
				if err != nil {
					fmt.Print(err)
				}
				return str
			case j.IsBool():
				bl, err := j.GetBool()
				if err != nil {
					fmt.Print(err)
				}
				return bl
			default:
				flt, err := j.GetFloat()
				if err != nil {
					fmt.Print(err)
				}
				return flt
			}
		}
	}
	return "No response from API. Check Config"
}
func get(ctx context.Context, pxrest *proffixrest.Client, endpoint string, params url.Values, format string, depth int, field string) (resp interface{}) {

	res, _, _, err := pxrest.Get(ctx, endpoint, params)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	return responseFormatter(res, format, depth, field)

}
func post(ctx context.Context, pxrest *proffixrest.Client, endpoint string, data string) (resp string) {

	_, header, _, _ := pxrest.Post(ctx, endpoint, data) // Ignore error, return empty on failure
	return proffixrest.ConvertLocationToID(header)
}
