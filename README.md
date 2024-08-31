[![GitHub release](https://img.shields.io/github/tag/pitwch/go-wrapper-proffix-restapi.svg?label=latest&colorB=orange)](https://github.com/pitwch/go-wrapper-proffix-restapi/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/pitwch/go-wrapper-proffix-restapi)](https://goreportcard.com/report/github.com/pitwch/go-wrapper-proffix-restapi)
[![Build Status]([https://travis-ci.com/pitwch/go-wrapper-proffix-restapi.svg?branch=master](https://app.travis-ci.com/pitwch/go-wrapper-proffix-restapi.svg))](https://app.travis-ci.com/github/pitwch/go-wrapper-proffix-restapi)
[![GoDoc](https://godoc.org/github.com/pitwch/go-wrapper-proffix-restapi/proffixrest?status.svg)](https://godoc.org/github.com/pitwch/go-wrapper-proffix-restapi/proffixrest)
[![codecov](https://codecov.io/gh/pitwch/go-wrapper-proffix-restapi/branch/master/graph/badge.svg)](https://codecov.io/gh/pitwch/go-wrapper-proffix-restapi)
[![GitHub license](https://img.shields.io/github/license/pitwch/go-wrapper-proffix-restapi.svg)](https://github.com/pitwch/go-wrapper-proffix-restapi/blob/master/LICENSE)

# Golang Wrapper für PROFFIX REST-API

Der zuverlässige Wrapper für die PROFFIX REST-API in Go.

![alt text](https://raw.githubusercontent.com/pitwch/go-wrapper-proffix-restapi/master/_assets/img/go-proffix-rest-wrapper.jpg "Golang Wrapper PROFFIX REST API")

**Übersicht**

- [Installation](#installation)
  - [Konfiguration](#konfiguration)
- [Optionen](#optionen)
  - [Methoden](#methoden)
  - [Spezielle Endpunkte](#spezielle-endpunkte)
- [CMD / Docker](#cmd-docker)
- [Weitere Beispiele](#weitere-beispiele)

### Installation

```bash
go get github.com/pitwch/go-wrapper-proffix-restapi/proffixrest
```

#### Konfiguration

Die Konfiguration wird dem Client mitgegeben:

| Konfiguration | Beispiel                 | Type          | Bemerkung                             |
|---------------|--------------------------|---------------|---------------------------------------|
| RestURL       | <https://myserver.ch:999>  | `string`      | URL der REST-API **ohne pxapi/v2/**   |
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
| Version          | v3                                                               | API-Version; Standard = v3                                     |
| APIPrefix        | /pxapi/                                                          | Prefix für die API; Standard = /pxapi/                         |
| LoginEndpoint    | PRO/Login                                                        | Endpunkt für Login; Standard = PRO/Login                       |
| UserAgent        | go-wrapper-proffix-restapi                                       | User Agent; Standard = go-wrapper-proffix-restapi              |
| Timeout          | 15                                                               | Timeout in Sekunden                                            |
| VerifySSL        | true                                                             | SSL prüfen                                                     |
| Batchsize        | 200                                                              | Batchgrösse für Batchrequests; Standard = 200                  |
| Log              | true                                                             | Aktiviert den Log für Debugging; Standard = false              |
| Client           | urlfetch.Client(ctx)                                             | HTTP-Client; Standard = http.DefaultClient                     |
| VolumeLicence    | false                                                            | Nutzt PROFFIX Volumenlizenzierung                              |

#### Methoden

| Parameter  | Typ           | Bemerkung                                                                                                |
|------------|---------------|----------------------------------------------------------------------------------------------------------|
| endpoint   | `string`      | Endpunkt der PROFFIX REST-API; z.B. ADR/Adresse,STU/Rapporte...                                          |
| data       | `interface{}` | Daten (werden automatisch in JSON konvertiert)                                                           |
| parameters | `url.Values`  | Parameter gemäss [PROFFIX REST API Docs](http://www.proffix.net/Portals/0/content/REST%20API/index.html) |

Folgende unterschiedlichen Methoden sind mit dem Wrapper möglich:

##### Get / Query

```golang

//Context
ctx := context.Background()

//Einfache Abfrage
     rc, _,statuscode, err := pxrest.Get(ctx,"ADR/Adresse/1", url.Values{})

//Abfrage mit Parametern
 param := url.Values{}
 param.Set("Filter", "Vorname@='Max'")

 rc,_, statuscode, err := pxrest.Get(ctx,"ADR/Adresse", param)

```

##### Put / Update

```golang
$pxrest =  new  Client(...)
var data map[string]interface{} = map[string]interface{}{
  "Name":   "Muster GmbH",
  "Ort":    "Zürich",
  "Zürich": "8000",
 }

 rc, header, statuscode, err := pxrest.Put(ctx,"ADR/Adresse", data)
```

##### Patch / Update

```golang
$pxrest =  new  Client(...)
var data map[string]interface{} = map[string]interface{}{
  "Ort":    "Zürich",
  "Zürich": "8000",
 }

 rc, header, statuscode, err := pxrest.Patch(ctx,"ADR/Adresse/1", data)
```

##### Post / Create

```golang
var data map[string]interface{} = map[string]interface{}{
  "Name":   "Muster GmbH",
  "Ort":    "Zürich",
  "Zürich": "8000",
 }

 //Query Endpoint ADR/Adresse with Headers
 rc, header, statuscode, err := pxrest.Post(ctx,"ADR/Adresse", data)
```

##### Delete / Delete

```golang
 //Delete Endpoint ADR/Adresse with Headers
 _, _, _, err = pxrest.Delete(ctx,"ADR/Adresse/2")
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
 rc, _, _, err := pxrest.Get(ctx,"ADR/Adresse/1", url.Values{})
```

Beispiel (Mit Header, Ohne Statuscode):

```golang
// Returns Header
 rc, header, _, err := pxrest.Post(ctx,"ADR/Adresse/1", data)

// Print Header->Location
 fmt.Print(header.Get("Location"))
 
```

**Hinweis:** Der Typ `io.ReadCloser` kann mittels Buffer Decode gelesen werden:

```golang
 rc, header, _, err := pxrest.Get(ctx,"ADR/Adresse/1", url.Values{})

 buf := new(bytes.Buffer)
 buf.ReadFrom(rc)
 resp := buf.String()
 fmt.Printf(resp, err)
 defer rc.Close()
```

##### Fehlerhandling / Error

Detaillierte Fehlerinformationen der REST-API können über den Fehlertyp PxError ermittelt werden.

```golang
// Query with Error
 rc, _, _, err := pxrest.Get(ctx,"ADR/Adresse/xxx", url.Values{})
    
    fmt.Print(err.(*PxError).Message)    // Message field
    fmt.Print(err.(*PxError).Status)     // Status field
    fmt.Print(err.(*PxError).Fields)     // Failed fields

```

Im Standard wird das Feld **Message** sowie fehlgeschlagene **Fields** ausgegeben.

#### Spezielle Endpunkte

##### Logout

Loggt den Client von der PROFFIX REST-API aus und gibt die Session / Lizenz damit wieder frei.

*Hinweis: Es wird automatisch die zuletzt verwendete PxSessionId für den Logout verwendet*

```golang

 statuscode, err := pxrest.Logout(ctx)

```

Der Wrapper führt den **Logout auch automatisch bei Fehlern** durch damit keine Lizenz geblockt wird.

##### Info

Ruft Infos vom Endpunkt **PRO/Info** ab.

*Hinweis: Dieser Endpunkt / Abfrage blockiert keine Lizenz*

```golang

 rc, err := pxrest.Info(ctx,"112a5a90fe28b23ed2c776562a7d1043957b5b79fad242b10141254b4de59028")

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
 rc, err := pxrest.Database(ctx,"112a5a90fe28b23ed2c776562a7d1043957b5b79fad242b10141254b4de59028")

 //Buffer decode for plain text response
 buf := new(bytes.Buffer)
 buf.ReadFrom(rc)
 resp := buf.String()

 fmt.Printf(resp, err)

 defer rc.Close()
  ```

##### ServiceLogin

Loggt den Client mit einer gegebenen PxSessionId ein (wird z.B. bei Erweiterungen verwendet)

```golang

 pxrest, err := pxrest.ServiceLogin(ctx,"7e024696-49e4-e54e-1a20-ffc0ba9a1374)

```

Der Wrapper führt den **Logout auch automatisch bei Fehlern** durch damit keine Lizenz geblockt wird.

##### PRO/Datei bzw. File Upload

Lädt Datei zum Endpunkt **PRO/Datei** hoch.

```golang

 fp := "C:/test.png"
 file, err := ioutil.ReadFile(fp)

 rc, headers, status, err := pxrest.File(ctx,"Demofile.png", file)

 //Buffer decode for plain text response
 buf := new(bytes.Buffer)
 buf.ReadFrom(rc)
 resp := buf.String()

 fmt.Printf(resp, err)

 defer rc.Close()

    // Get temporary filename from Headers
    id := proffixrest.ConvertLocationToID(headers)

```

##### PRO/Datei bzw. File Download

Liest eine Datei aufgrund der DateiNr

```golang

rc, filename,_,_,err := px.GetFile(ctx,"9PnrK4XiFzyrTdDK9HHLE80Kk2lARKS5Ppn2WNrOPqGHzCREYtidmRhUKXq-82uSvDerMckkcB3srYn9TsKvk0C1P9ca3deIeMzp5jiaeIwT4wgnH0qarwdny17-8eGq",nil)
px.WriteFile(filename,rc)

```

##### GET Batch

Gibt sämtliche Ergebnisse aus und iteriert selbständig über die kompletten Ergebnisse der REST-API.
Ideal wenn z.B. nicht klar ist wieviele Ergebnisse vorliegen aber trotzdem alle gebraucht werden.

*Hinweis: Das mögliche Total der Ergebnisse wird aus dem Header `FilteredCount` ermittelt*

```golang
    //Set Params
     params := url.Values{}
     params.Set("Fields", "AdressNr,Name,Ort,Plz")

     //GET Batch on ADR/Adressen with Batchsize = 35
     rc, total, err := pxrest.GetBatch(ctx,"ADR/Adresse", params, 35)

     //Returns all possible results in `rc`, the complete amount of results as int in `total` and
     //possible errors in `err`


```

##### Sync Batch

Synchronisiert Daten im Batch Modus.
Damit lassen sich **beliebig viele Einträge eines Endpunktes (z.B: 100 Artikel) synchronisieren** - der Wrapper wählt **automatisch die passende Methode** aus und aktualisiert den JSON-Body (abhängig vom Vorhandensein eines Keyfields und dessen Ergebnis einer GET-Abfrage)

```golang
    
    //First entry has AdressNr and exist -> PUT, second has no AdressNr -> POST
   jsonExample := `[{
                    "AdressNr": 276,
                    "Name": "EYX AG",
                    "Vorname": null,
                    "Strasse": "Zollstrasse 2",
                    "PLZ": "8000",
                    "Ort": "Zürich"
                   }, {
                    "Name": "Muster",
                    "Vorname": "Hans",
                    "Strasse": "Zürcherstrasse 2",
                    "PLZ": "8000",
                    "Ort": "Zürich"
                   }]`
                     


   created,updated,failed,errors, total, err := pxrest.SyncBatch(ctx, "ADR/Adresse","AdressNr", true,[]byte(jsonExample)) 
 //Results
 //      created -> [333]
 //      updated -> [276]
 //      failed -> []
    //      errors -> []
    //      total -> 2
    //      err -> nil

```

Der Batch läuft von Anfang bis Ende durch - die Keyfelder der Ergebnisse werden in den jeweiligen Variablen `updated`,
`created`,`failed`,`errors` ausgegeben.

*Hinweis: Der Parameter **Keyfield** wird genutzt um je nach Methode das Schlüsselfeld im URL - Slug automatisch
anzupassen. Der Parameter **removeKeyfield** wird verwendet um das Keyfield aus dem Body zu entfernen (z.B: bei ADR/Adressen)*

##### GET List

Gibt direkt die Liste der PROFFIX REST API aus (ohne Umwege)

```golang
    //Get File
 file, headers, status, err := pxrest.GetList(ctx, "ADR_Adress-Etiketten Kontakte.labx", nil)

    //Write File
    px.WriteFile("C://test//test.pdf",file)  //Hilfsfunktion um Dateien zu schreiben

```

*Hinweis: Der Dateityp (zurzeit nur PDF) kann über den Header `File-Type` ermittelt werden*

##### CheckApi

Prüft, ob die API funktioniert (Basis-Test)

```golang
    //CheckApi
 err := pxrest.CheckApi(ctx)

```

*Hinweis: Die Test - Funktion loggt sich ein, prüft ob eine PxSessionId vorhanden ist und ruft eine Adresse ab.

#### Hilfsfunktionen

##### ConvertPXTimeToTime

Konvertiert einen Zeitstempel der PROFFIX REST-API in time.Time

```golang

t := ConvertPXTimeToTime('2004-04-11 00:00:00')

```

##### ConvertTimeToPXTime

Konvertiert einen time.Time in einen Zeitstempel der PROFFIX REST-API

```golang

t := ConvertTimeToPXTime(time.Now())

```

##### ConvertLocationToID

Extrahiert die ID aus dem Header Location der PROFFIX REST-API

```golang

id := ConvertLocationToID(header)

```

##### WriteFile

Schreibt eine Datei

```golang

px.WriteFile("C://test//test.pdf",file)

```

### CMD / Docker

todo

### Weitere Beispiele

Im Ordner [/examples](https://github.com/pitwch/go-wrapper-proffix-restapi/tree/master/_examples) finden sich weitere,
auskommentierte Beispiele.
