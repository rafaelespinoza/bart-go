package bart

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestTripPlans(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqURL, err := url.ParseRequestURI(r.RequestURI)
		if err != nil {
			panic(err)
		}
		if reqURL.Path != "/sched.aspx" {
			t.Errorf("Wrong request path. Got %s, expected %s", reqURL.Path, "/sched.aspx")
			return
		}
		query := reqURL.Query()
		var a, b int
		var aOK, bOK bool
		var params []string
		if params, aOK = query["a"]; aOK {
			a, _ = strconv.Atoi(params[0])
		}
		if params, bOK := query["b"]; bOK {
			b, _ = strconv.Atoi(params[0])
		}

		if (aOK && bOK) && (b == 0 && a == 0) || (b == 0 && a == 1) {
			// failure expected in this case
			fmt.Fprint(w, `
{
	"root":{
		"schedule": {
			"request": {
				"trip": {
					"@origTimeMin": "07:37 PM",
					"@destTimeMin": "07:44 PM"
				}
			}
		}
	}
}`)
		} else {
			fmt.Fprint(w, `
{
	"root":{
		"schedule": {
			"request": {
				"trip": [
					{ "@origTimeMin": "07:37 PM", "@destTimeMin": "07:44 PM" },
					{ "@origTimeMin": "07:41 PM", "@destTimeMin": "07:48 PM" }
				]
			}
		}
	}
}`)
		}
	}))

	client := NewClient(nil)
	client.conf.baseURL = server.URL
	defer func() {
		server.Close()
	}()

	testTripsResponse := func(t *testing.T, res TripsResponse, err error) {
		t.Helper()
		if err != nil {
			t.Error(err.Error())
			return
		}

		list := res.Root.Data.Request.List
		if len(list) < 1 {
			t.Error("no results at Root.Data.Request.List")
			return
		}

		for i := range list {
			if list[i].OrigTimeMin == "" {
				t.Error("expected non-empty value for OrigTimeMin")
			}
			if list[i].DestTimeMin == "" {
				t.Error("expected non-empty value for DestTimeMin")
			}
		}
	}

	t.Run("Arrivals", func(t *testing.T) {
		var params TripParams
		var res TripsResponse
		var err error

		params = TripParams{Orig: "woak", Dest: "embr"}
		res, err = client.RequestArrivals(params)
		testTripsResponse(t, res, err)

		params = TripParams{Orig: "woak", Dest: "embr", Before: 0, After: 1}
		res, err = client.RequestArrivals(params)
		testTripsResponse(t, res, err)

		params = TripParams{Orig: "woak", Dest: "embr", Before: 0, After: 4}
		res, err = client.RequestArrivals(params)
		testTripsResponse(t, res, err)
	})

	t.Run("Departures", func(t *testing.T) {
		var params TripParams
		var res TripsResponse
		var err error

		params = TripParams{Orig: "woak", Dest: "embr"}
		res, err = client.RequestDepartures(params)
		testTripsResponse(t, res, err)

		params = TripParams{Orig: "woak", Dest: "embr", Before: 0, After: 1}
		res, err = client.RequestDepartures(params)
		testTripsResponse(t, res, err)

		params = TripParams{Orig: "woak", Dest: "embr", Before: 0, After: 4}
		res, err = client.RequestDepartures(params)
		testTripsResponse(t, res, err)
	})
}

func TestSpecialSchedules(t *testing.T) {
	t.Run("empty data", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqURL, err := url.ParseRequestURI(r.RequestURI)
			if err != nil {
				panic(err)
			}
			if reqURL.Path != "/sched.aspx" {
				t.Errorf("Wrong request path. Got %s, expected %s", reqURL.Path, "/sched.aspx")
				return
			}
			cmd := reqURL.Query().Get("cmd")
			switch cmd {
			case "special":
				fmt.Fprintf(w, `
{
   "?xml":{
      "@version":"1.0",
      "@encoding":"utf-8"
   },
   "root":{
      "uri":{
         "#cdata-section":"http://api.bart.gov/api/sched.aspx?cmd=special&l=1&json=y"
      },
      "special_schedules":"",
      "message":""
   }
}
			`)
			default:
				err := fmt.Errorf("Unknown query string value %s", cmd)
				panic(err)
			}
		}))
		client := NewClient(nil)
		client.conf.baseURL = server.URL
		defer func() {
			server.Close()
		}()
		_, err := client.RequestSpecialSchedules()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		// should be able to parse w/o error
	})

	t.Run("non-empty data", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqURL, err := url.ParseRequestURI(r.RequestURI)
			if err != nil {
				panic(err)
			}
			if reqURL.Path != "/sched.aspx" {
				t.Errorf("Wrong request path. Got %s, expected %s", reqURL.Path, "/sched.aspx")
				return
			}
			cmd := reqURL.Query().Get("cmd")
			switch cmd {
			case "special":
				fmt.Fprintf(w, `
{
    "root": {
        "uri": {
            "#cdata-section": "http://api.bart.gov/api/sched.aspx?cmd=special&l=1&json=y"
        },
        "special_schedules": {
            "special_schedule": [
                {
                    "start_date": "07/01/2017",
                    "end_date": "07/02/2017",
                    "start_time": "",
                    "end_time": "",
                    "text": {
                        "#cdata-section": "Expect delays of 20 to 40 minutes because of bus bridge between Fruitvale & 19th St."
                    },
                    "link": {
                        "#cdata-section": "http://www.bart.gov/news/articles/2017/news20170302-1"
                    },
                    "orig": "",
                    "dest": "",
                    "day_of_week": "",
                    "routes_affected": "ROUTE 3, ROUTE 4, ROUTE 5, ROUTE 6, ROUTE 11, ROUTE 12, ROUTE 19, ROUTE 20"
                }
            ],
            "message": ""
        }
    }
}`)
			default:
				err := fmt.Errorf("Unknown query string value %s", cmd)
				panic(err)
			}
		}))
		client := NewClient(nil)
		client.conf.baseURL = server.URL
		defer func() {
			server.Close()
		}()
		out, err := client.RequestSpecialSchedules()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		val := out.Root.Data
		if len(val.List) < 1 {
			t.Error("Expected non-empty data")
			return
		}
		if len(val.List[0].StartDate) < 1 {
			t.Errorf("Expected value for StartDate")
		}
		if len(val.List[0].EndDate) < 1 {
			t.Errorf("Expected value for EndDate")
		}
		if len(val.List[0].Text.Value) < 1 {
			t.Errorf("Expected value for Text.Value")
		}
		if len(val.List[0].RoutesAffected) < 1 {
			t.Errorf("Expected value for RoutesAffected")
		}
	})
}
