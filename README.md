[![Go Report Card](https://goreportcard.com/badge/github.com/pitwch/go-wrapper-proffix-restapi)](https://goreportcard.com/report/github.com/pitwch/go-wrapper-proffix-restapi)
[![Build Status](https://travis-ci.org/pitwch/go-wrapper-proffix-restapi.svg?branch=master)](https://travis-ci.org/pitwch/go-wrapper-proffix-restapi)
[![GoDoc](https://godoc.org/github.com/pitwch/go-wrapper-proffix-restapi/proffixrest?status.svg)](https://godoc.org/github.com/pitwch/go-wrapper-proffix-restapi/proffixrest)
[![codecov](https://codecov.io/gh/pitwch/go-wrapper-proffix-restapi/branch/master/graph/badge.svg)](https://codecov.io/gh/pitwch/go-wrapper-proffix-restapi)


# Golang Wrapper für PROFFIX REST-API

Der zuverlässige Wrapper für die PROFFIX REST-API in Go.

![alt text](https://raw.githubusercontent.com/pitwch/go-wrapper-proffix-restapi/master/assets/img/go-proffix-rest-wrapper.jpg "Golang Wrapper PROFFIX REST API")


**Übersicht**

- [Installation](#installation)
  - [Konfiguration](#konfiguration)
- [Optionen](#optionen)
  - [Methoden](#methoden)
  - [Spezielle Endpunkte](#spezielle-endpunkte)
- [Weitere Beispiele](#weitere-beispiele)


### Installation

```bash
$ go get github.com/pitwch/go-wrapper-proffix-restapi/proffixrest
```


#### Konfiguration

Die Konfiguration wird dem Client mitgegeben:

| Konfiguration | Beispiel                 | Type          | Bemerkung                             |
|---------------|--------------------------|---------------|---------------------------------------|
| Context       | ctx                      | `context`     | Context                               |
| RestURL       | https://myserver.ch:999  | `string`      | URL der REST-API **ohne pxapi/v2/**   |
| apiDatabase   | DEMO                     | `string`      | Name der Datenbank                    |
| apiUser       | USR                      | `string`      | Names des Benutzers                   |
| apiPassword   | b62cce2fe18f7a156a9c...  | `string`      | SHA256-Hash des Benutzerpasswortes    |
| apiModule     | []string{"ADR", "FIB"}   | `[]string`    | Benötigte Module (mit Komma getrennt) |
| options       | &px.Options{Timeout: 30} | `*px.Options` | Optionen (Details unter Optionen)     |

Beispiel:

```golang
import (
  px "github.com/pitwch/go-wrapper-proffix-restapi/proffixrest"
)

//Set Context
ctx := context.Background()

var pxrest, err = px.NewClient(
    ctx,
	"https://myserver.ch:999",
	"USR",
	"b62cce2fe18f7a156a9c719c57bebf0478a3d50f0d7bd18d9e8a40be2e663017",
	"DEMO",
	[]string{"ADR", "FIB"},
	&px.Options{},
)
```


### Optionen

Optionen sind **fakultativ** und werden in der Regel nicht benötigt:

| Option           | Beispiel                                                         | Bemerkung                                                      |
|------------------|------------------------------------------------------------------|----------------------------------------------------------------|
| Key              | 112a5a90fe28b...242b10141254b4de59028                            | API-Key als SHA256 - Hash (kann auch direkt mitgegeben werden) |
| Version          | v2                                                               | API-Version; Standard = v2                                     |
| APIPrefix        | /pxapi/                                                          | Prefix für die API; Standard = /pxapi/                         |
| LoginEndpoint    | PRO/Login                                                        | Endpunkt für Login; Standard = PRO/Login                       |
| UserAgent        | go-wrapper-proffix-restapi                                       | User Agent; Standard = go-wrapper-proffix-restapi              |
| Timeout          | 15                                                               | Timeout in Sekunden                                            |
| VerifySSL        | true                                                             | SSL prüfen                                                     |
| Batchsize        | 200                                                              | Batchgrösse für Batchrequests; Standard = 200                  |
| Log              | true                                                             | Aktiviert den Log für Debugging; Standard = false              |
| Client           | urlfetch.Client(ctx)                                             | HTTP-Client; Standard = http.DefaultClient                     |




#### Methoden


| Parameter  | Typ           | Bemerkung                                                                                                |
|------------|---------------|----------------------------------------------------------------------------------------------------------|
| endpoint   | `string`      | Endpunkt der PROFFIX REST-API; z.B. ADR/Adresse,STU/Rapporte...                                          |
| data       | `interface{}` | Daten (werden automatisch in JSON konvertiert)                                                           |
| parameters | `url.Values`  | Parameter gemäss [PROFFIX REST API Docs](http://www.proffix.net/Portals/0/content/REST%20API/index.html) |




Folgende unterschiedlichen Methoden sind mit dem Wrapper möglich:



##### Get / Query

```golang
//Einfache Abfrage
     rc, _,statuscode, err := pxrest.Get("ADR/Adresse/1", url.Values{})

//Abfrage mit Parametern
	param := url.Values{}
	param.Set("Filter", "Vorname@='Max'")

	rc,_, statuscode, err := pxrest.Get("ADR/Adresse", param)

```


##### Put / Update

```golang
$pxrest =  new  Client(...)
var data map[string]interface{} = map[string]interface{}{
		"Name":   "Muster GmbH",
		"Ort":    "Zürich",
		"Zürich": "8000",
	}

	rc, header, statuscode, err := pxrest.Put("ADR/Adresse", data)
```

##### Post / Create

```golang
var data map[string]interface{} = map[string]interface{}{
		"Name":   "Muster GmbH",
		"Ort":    "Zürich",
		"Zürich": "8000",
	}

	//Query Endpoint ADR/Adresse with Headers
	rc, header, statuscode, err := pxrest.Post("ADR/Adresse", data)
```


##### Delete / Delete

```golang
	//Delete Endpoint ADR/Adresse with Headers
	_, _, _, err = pxrest.Delete("ADR/Adresse/2")
```

##### Response / Antwort

Alle Methoden geben `io.ReadCloser`, `http.Header`,  `int` sowie `error` zurück.

| Rückgabetyp     | Bemerkung                                      |
|-----------------|------------------------------------------------|
| `io.ReadCloser` | Daten                                          |
| `http.Header`   | Header                                         |
| `int`           | HTTP-Status Code                               |
| `error`         | Fehler; Standardrückgabe ohne Fehler ist `nil` |

Beispiel (Kein Header, Ohne Statuscode):

```golang
// Returns no Header (_)
	rc, _, _, err := pxrest.Get("ADR/Adresse/1", url.Values{})
```

Beispiel (Mit Header, Ohne Statuscode):

```golang
// Returns Header
	rc, header, _, err := pxrest.Post("ADR/Adresse/1", data)

// Print Header->Location
	fmt.Print(header.Get("Location"))
	
```

**Hinweis:** Der Typ `io.ReadCloser` kann mittels Buffer Decode gelesen werden:

```golang
	rc, header, _, err := pxrest.Get("ADR/Adresse/1", url.Values{})

	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()
	fmt.Printf(resp, err)
	defer rc.Close()
```

#### Spezielle Endpunkte

##### Logout

Loggt den Client von der PROFFIX REST-API aus und gibt die Session / Lizenz damit wieder frei.

*Hinweis: Es wird automatisch die zuletzt verwendete PxSessionId für den Logout verwendet*


```golang

	statuscode, err := pxrest.Logout()

```

Der Wrapper führt den **Logout auch automatisch bei Fehlern** durch damit keine Lizenz geblockt wird.

##### Info

Ruft Infos vom Endpunkt **PRO/Info** ab.

*Hinweis: Dieser Endpunkt / Abfrage blockiert keine Lizenz*

```golang

	rc, err := pxrest.Info("112a5a90fe28b23ed2c776562a7d1043957b5b79fad242b10141254b4de59028")

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	fmt.Printf(resp, err)

	defer rc.Close()
```

##### Datenbank

Ruft Infos vom Endpunkt **PRO/Datenbank** ab.

*Hinweis: Dieser Endpunkt / Abfrage blockiert keine Lizenz*

```golang
	rc, err := pxrest.Database("112a5a90fe28b23ed2c776562a7d1043957b5b79fad242b10141254b4de59028")

	//Buffer decode for plain text response
	buf := new(bytes.Buffer)
	buf.ReadFrom(rc)
	resp := buf.String()

	fmt.Printf(resp, err)

	defer rc.Close()
  ```

##### GET Batch

Gibt sämtliche Ergebnisse aus und iteriert selbständig über die kompletten Ergebnisse der REST-API.

*Hinweis: Das mögliche Total der Ergebnisse wird aus dem Header `FilteredCount` ermittelt*

```golang
    //Set Params
     params := url.Values{}
     params.Set("Fields", "AdressNr,Name,Ort,Plz")

     //GET Batch on ADR/Adressen with Batchsize = 35
     rc, total, err := pxrest.GetBatch("ADR/Adresse", params, 35)

     //Returns all possible results in `rc`, the complete amount of results as int in `total` and
     //possible errors in `err`


```


### Weitere Beispiele

Im Ordner [/examples](https://github.com/pitwch/go-wrapper-proffix-restapi/tree/master/examples) finden sich weitere,
auskommentierte Beispiele.
