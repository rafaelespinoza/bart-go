# bart

[![GoDoc](https://godoc.org/github.com/rafaelespinoza/bart-go?status.svg)](https://godoc.org/github.com/rafaelespinoza/bart-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/rafaelespinoza/bart-go)](https://goreportcard.com/report/github.com/rafaelespinoza/bart-go)

BART API client in golang.
Supports Advisories, Real-Time Estimates, Routes, Schedules, Stations APIs.
Outputs JSON.

`go get -u github.com/rafaelespinoza/bart-go`

-   [Official BART API docs](https://api.bart.gov/docs/overview/index.aspx)
-   [Package docs](https://godoc.org/github.com/rafaelespinoza/bart-go)


See the examples file for general usage.

---

The response schema from the BART API is a little irregular and this package makes every attempt to
make field names more consistent across APIs. For example, a response from `/stn.aspx?cmd=stns&json=y` gives you:

```json
{
    "root": {
        "stations": {
            "station": [ { } ]
        }
    }
}
```

While the resulting go struct will be in the shape:

```go
    Root struct {
        Data struct {
            List []struct { }
        }
    }
```

#### station names

There are several methods that require an orig or dest value to be the name of the station. Valid
values for orig, dest inputs are 4-letter abbreviations for the station name. [Here is a full list
of station abbreviations](https://api.bart.gov/docs/overview/abbrev.aspx). The methods in this
package will accept those values as upper, lower or mixed case. If passed an invalid value, an error
is returned instead of performing the request.

#### available schedules, schedule numbers

The BART API does let you query for results based on past or future schedules, but this package
elects to use the current schedule only.
