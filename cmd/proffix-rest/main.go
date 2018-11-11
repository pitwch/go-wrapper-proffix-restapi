package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pitwch/go-wrapper-proffix-restapi/proffixrest"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
)
var SRVERSION string

func main() {
	resturl := flag.String("url", os.Getenv("PXR_URL"), "URL Rest-API")
	user := flag.String("user", os.Getenv("PXR_USER"), "Benutzer")
	password := flag.String("password", os.Getenv("PXR_PASSWORD"), "Passwort")
	database := flag.String("database", os.Getenv("PXR_DATABASE"), "Datenbank")
	module := flag.String("module", os.Getenv("PXR_MODULE"), "Module")
	apikey := flag.String("apikey", os.Getenv("PXR_APIKEY"), "API Passwort / Key")
	method := flag.String("method", "", "Methode")
	endpoint := flag.String("endpoint", "", "Endpunkt")
	data := flag.String("data", "", "Daten")

	limit := flag.String("limit", "1", "Limit")
	format := flag.String("format", "json", "Format")
	field := flag.String("field", "", "Feld")

	showVersion := flag.Bool("version", false, "outputs the semantic-release version")
	//updateFile := flag.String("update", "", "updates the version of a certain file")

	flag.Parse()

	if *showVersion {
		fmt.Printf("go-wrapper-proffix-restapi v%s", SRVERSION)
		return
	}

	param := url.Values{}
	param.Add("limit", *limit)

	ctx := context.Background()
	pxrest, err := proffixrest.NewClient(
		*resturl,
		*user,
		*password,
		*database,
		[]string{*module},
		&proffixrest.Options{Key: *apikey},
	)

	if err != nil {
		fmt.Printf("Fehler: %v", err)
	}

	switch strings.ToLower(*method) {
	case "get":
		fmt.Print(get(*pxrest, ctx, *endpoint, param, *format, *field))
	case "post":
		fmt.Print(post(*pxrest, ctx, *endpoint, *data))
	}

}
func responseFormatter(closer io.ReadCloser, format string, field string) interface{} {
	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(closer)
	respStr := buf.String()
	respByte := buf.Bytes()
	defer closer.Close()
	switch strings.ToLower(format) {
	case "json":
		return respStr
	case "string":
		m := make(map[string]interface{})
		err := json.Unmarshal(respByte, &m)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(m[field])
	}
	return ""
}
func get(pxrest proffixrest.Client, ctx context.Context, endpoint string, params url.Values, format string, field string) (resp interface{}) {

	res, _, _, _ := pxrest.Get(ctx, endpoint, params)

	return responseFormatter(res, format, field)

}
func post(pxrest proffixrest.Client, ctx context.Context, endpoint string, data string) (resp string) {

	_, header, _, _ := pxrest.Post(ctx, endpoint, data)
	return proffixrest.ConvertLocationToID(header)
}
