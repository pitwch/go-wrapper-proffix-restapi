### Golang Wrapper für PROFFIX REST-API



### Installation

```bash
$ go get github.com/pitw/go-wrapper-proffix-restapi/proffixrest
```


#### Konfiguration

Die Konfiguration wird dem Client mitgegeben:

| Konfiguration | Beispiel                                          | Bemerkung                             |
|---------------|---------------------------------------------------|---------------------------------------|
| RestURL       | https://myserver.ch:999                           | URL der REST-API **ohne pxapi/v2/**   |
| apiDatabase   | DEMO                                              | Name der Datenbank                    |
| apiUser       | USR                                               | Names des Benutzers                   |
| apiPassword   | b62cce2fe18f7a156a9c...0f0d7bd18d9e8a40be2e663017 | SHA256-Hash des Benutzerpasswortes    |
| apiModule     | ADR,STU                                           | Benötigte Module (mit Komma getrennt) |
| options       | array('key'=>'112a5a90...59028')                  | Optionen (Details unter Optionen)     |
Beispiel:

```golang
import (
  px "github.com/pitw/go-wrapper-proffix-restapi/proffixrest"
)

var pxrest, err = px.NewClient(
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
| UserAgent        | php-wrapper-proffix-restapi                                      | User Agent; Standard = php-wrapper-proffix-restapi             |
| Timeout          | 15                                                               | Timeout in Sekunden                                            |
| VerifySSL        | true                                                             | SSL prüfen                                                     |


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
     rc, _, err := pxrest.Get("ADR/Adresse/1", url.Values{})

/Abfrage mit Parametern
	param := url.Values{}
	param.Set("Filter", "Vorname@='Max'")

	rc,_, err := pxrest.Get("ADR/Adresse", param)

```


##### Put / Update

```golang
$pxrest =  new  Client(...)
var data map[string]interface{} = map[string]interface{}{
		"Name":   "Muster GmbH",
		"Ort":    "Zürich",
		"Zürich": "8000",
	}

	rc, header, err := pxrest.Put("ADR/Adresse", data)
```

##### Post / Create

```golang
var data map[string]interface{} = map[string]interface{}{
		"Name":   "Muster GmbH",
		"Ort":    "Zürich",
		"Zürich": "8000",
	}

	//Query Endpoint ADR/Adresse with Headers
	rc, header, err := pxrest.Post("ADR/Adresse", data)
```


##### Response / Antwort


todo


#### Spezielle Endpunkte


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
### Ausnahmen / Spezialfälle

* Endpunkte welche Leerschläge enthalten können (z.B. LAG/Artikel/PC 7/Bestand) müssen mit rawurlencode() genutzt werden

### Weitere Beispiele

Im Ordner [/examples](https://github.com/pitwch/go-wrapper-proffix-restapi/tree/master/examples) finden sich weitere,
auskommentierte Beispiele.