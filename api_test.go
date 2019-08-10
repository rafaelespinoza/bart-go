package bart

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type stubAPI struct {
	conf *Config
}

func (a *stubAPI) clientConf() *Config {
	if a != nil && a.conf != nil {
		return a.conf
	}
	return defaultClientConf
}

func TestRequestAPI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqURL, err := url.ParseRequestURI(r.RequestURI)
		if err != nil {
			panic(err)
		}
		switch reqURL.Path {
		case "/route.aspx":
			fmt.Fprintf(w, `<?xml version="1.0" encoding="utf-8"?>
<root>
	<message>
		<error>
			<text>Invalid cmd</text>
			<details>The cmd parameter (%s) is missing or invalid. Please correct the error and try again.</details>
		</error>
	</message>
</root>`, reqURL.Query().Get("cmd"))
		default:
			err := fmt.Errorf("unknown request path %s", reqURL.Path)
			panic(err)
		}
	}))
	client := stubAPI{
		conf: &Config{
			Key:  Key,
			HTTP: &http.Client{},
		},
	}
	client.conf.baseURL = server.URL
	defer server.Close()

	var out interface{}
	actual := requestAPI(&client, "/route.aspx", "foo", nil, &out)
	expectedError := "error: Invalid cmd. The cmd parameter (foo) is missing or invalid. Please correct the error and try again."
	if actual.Error() != expectedError {
		t.Errorf("Error formatted incorrectly. Got %q, expected %q", actual.Error(), expectedError)
	}
}

type transporter struct{}

func (tr *transporter) RoundTrip(req *http.Request) (*http.Response, error) {
	return http.DefaultTransport.RoundTrip(req)
}

func TestClient(t *testing.T) {
	var expectedURI string
	apiClient := NewClient(&Config{
		Key:  "FOO-BAR",
		HTTP: &http.Client{Transport: &transporter{}},
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualReqURI := r.URL.String()
		if actualReqURI != expectedURI {
			t.Errorf("wrong URI. Got %q, expected %q", actualReqURI, expectedURI)
		}
		fmt.Fprint(w, `{"root": {"uri": { "#cdata-section": "ok usa" }}}`)
	}))

	apiClient.conf.baseURL = server.URL
	defaultClientConf.baseURL = server.URL
	defer func() {
		defaultClientConf.baseURL = baseURL
		server.Close()
	}()

	tests := []struct {
		client *Client
		apiKey string
	}{
		{client: apiClient, apiKey: "FOO-BAR"},

		// Previously documented examples, which used `new`, should work.
		{client: new(Client), apiKey: Key},
	}

	for _, test := range tests {
		client := test.client
		key := test.apiKey

		// AdvisoriesAPI
		expectedURI = fmt.Sprintf("/bsa.aspx?cmd=bsa&json=y&key=%s", key)
		client.RequestBSA()
		expectedURI = fmt.Sprintf("/bsa.aspx?cmd=elev&json=y&key=%s", key)
		client.RequestElevator()
		expectedURI = fmt.Sprintf("/bsa.aspx?cmd=count&json=y&key=%s", key)
		client.RequestTrainCount()

		// EstimatesAPI
		expectedURI = fmt.Sprintf("/etd.aspx?cmd=etd&dir=s&json=y&key=%s&orig=mcar&plat=4", key)
		client.RequestETD("mcar", "4", "s")
		expectedURI = fmt.Sprintf("/etd.aspx?cmd=etd&dir=n&json=y&key=%s&orig=civc", key)
		client.RequestEstimate(EstimateParams{Orig: "civc", Dir: "n"})

		// RoutesAPI
		expectedURI = fmt.Sprintf("/route.aspx?cmd=routes&json=y&key=%s", key)
		client.RequestRoutes("")
		expectedURI = fmt.Sprintf("/route.aspx?cmd=routeinfo&json=y&key=%s&route=all", key)
		client.RequestRoutesInfo("")

		// SchedulesAPI
		expectedURI = fmt.Sprintf("/sched.aspx?a=3&b=3&cmd=arrive&dest=cols&json=y&key=%s&l=1&orig=embr&time=6%s30pm", key, "%3A")
		client.RequestArrivals(TripParams{
			Orig:   "embr",
			Dest:   "cols",
			Time:   "6:30pm",
			Date:   "",
			Before: 3,
			After:  3,
			Legend: true,
		})
		expectedURI = fmt.Sprintf("/sched.aspx?cmd=scheds&json=y&key=%s", key)
		client.RequestAvailableSchedules()
		expectedURI = fmt.Sprintf("/sched.aspx?a=2&b=2&cmd=depart&dest=cast&json=y&key=%s&l=1&orig=sfia", key)
		client.RequestDepartures(TripParams{
			Orig:   "sfia",
			Dest:   "cast",
			Time:   "",
			Date:   "",
			Before: 2,
			After:  2,
			Legend: true,
		})
		expectedURI = fmt.Sprintf("/sched.aspx?cmd=holiday&json=y&key=%s", key)
		client.RequestHolidaySchedules()
		expectedURI = fmt.Sprintf("/sched.aspx?cmd=routesched&date=sa&json=y&key=%s&l=1&route=12", key)
		client.RequestRouteSchedules(12, "sa", "", true)
		expectedURI = fmt.Sprintf("/sched.aspx?cmd=special&json=y&key=%s", key)
		client.RequestSpecialSchedules()
		expectedURI = fmt.Sprintf("/sched.aspx?cmd=stnsched&json=y&key=%s&orig=mcar", key)
		client.RequestStationSchedules("mcar", "")
		expectedURI = fmt.Sprintf("/sched.aspx?cmd=stnsched&date=8%s14%s2018&json=y&key=%s&orig=glen", "%2F", "%2F", key)
		client.RequestStationSchedules("glen", "8/14/2018")

		// StationsAPI
		expectedURI = fmt.Sprintf("/stn.aspx?cmd=stnaccess&json=y&key=%s&orig=phil", key)
		client.RequestStationAccess("phil")
		expectedURI = fmt.Sprintf("/stn.aspx?cmd=stninfo&json=y&key=%s&orig=mont", key)
		client.RequestStationInfo("mont")
		expectedURI = fmt.Sprintf("/stn.aspx?cmd=stns&json=y&key=%s", key)
		client.RequestStations()
	}
}
