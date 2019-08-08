// Package bart is a wrapper for the BART API. It works with JSON, which is
// still in beta at the time of this writing. See the official BART docs for
// information https://api.bart.gov/docs/overview/index.aspx
//
// The API request example functions are written with vague regard to the
// output. However, to execute the examples with go test, we must specify some
// kind of output. Making real requests to the BART API in tests (as opposed to
// stubbing responses) is meant to expose unexpected errors or panics, which
// should help make response handling better in the long run.
package bart

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	// Key is default API Key that BART gives to all developers. If you have
	// registered your own key, then you should use the NewClient function.
	Key     = "MW9S-E7SL-26DU-VV8V"
	baseURL = "https://api.bart.gov/api"
)

var defaultClientConf = &Config{
	Key:     Key,
	HTTP:    &http.Client{},
	baseURL: baseURL,
}

// A Config is a collection of named parameters for a Client.
type Config struct {
	Key     string
	HTTP    *http.Client
	baseURL string
}

// Client gives you easy access to several BART API endpoints. See examples for
// general usage.
type Client struct {
	*AdvisoriesAPI
	*EstimatesAPI
	*RoutesAPI
	*SchedulesAPI
	*StationsAPI

	conf *Config
}

// NewClient allows you to set up a Client. Pass in nil if you want to use the
// default settings. If you have registered our own API key, then specify
// conf.Key. If conf.Key is empty, then the default API key is used. If
// conf.HTTP is empty then the http client is an empty *http.Client from the
// standard library.
func NewClient(conf *Config) *Client {
	if conf == nil {
		conf = &Config{}
	}
	if len(conf.Key) == 0 {
		conf.Key = Key
	}
	if conf.HTTP == nil {
		conf.HTTP = &http.Client{}
	}
	conf.baseURL = baseURL
	return &Client{
		conf:          conf,
		AdvisoriesAPI: &AdvisoriesAPI{conf},
		EstimatesAPI:  &EstimatesAPI{conf},
		SchedulesAPI:  &SchedulesAPI{conf},
		StationsAPI:   &StationsAPI{conf},
		RoutesAPI:     &RoutesAPI{conf},
	}
}

// ResponseMetaData is contains some data about the response. Not all of the
// fields are filled by every API endpoint.
type ResponseMetaData struct {
	URI     CDATASection
	Date    string      `json:",omitempty"`
	Time    string      `json:",omitempty"`
	Message interface{} `json:",omitempty"`
}

// CDATASection is merely a helper for unmarshaling certain fields. The original
// BART API has long returned XML instead of JSON, and its presence is an
// artifact of BART's output conversion. This type is meant for internal use and
// is only exported for documentation purposes since it shows up in so many
// other type definitions in this package.
type CDATASection struct {
	Value string `json:"#cdata-section"`
}

type configuredClient interface {
	clientConf() *Config
}

func requestAPI(cc configuredClient, route, cmd string, params *url.Values, out interface{}) (err error) {
	var res *http.Response

	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	conf := cc.clientConf()
	if params == nil {
		params = &url.Values{}
	}
	params.Set("cmd", cmd)
	params.Set("json", "y")
	params.Set("key", conf.Key)
	uri := conf.baseURL + route + "?" + params.Encode()

	res, err = conf.HTTP.Get(uri)
	if err != nil {
		return err
	}

	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(raw, out)
	if _, ok := (err).(*json.SyntaxError); ok {
		// Handle any errors from the BART API that are formatted as XML. As of
		// this writing, the JSON API format is in beta and they never got
		// around to converting errors to JSON.
		var xmlMessage struct {
			Text    string `xml:"message>error>text"`
			Details string `xml:"message>error>details"`
		}
		xmlParseErr := xml.Unmarshal(raw, &xmlMessage)
		if xmlParseErr != nil {
			err = xmlParseErr
			return
		}
		err = fmt.Errorf("error: %s. %s", xmlMessage.Text, xmlMessage.Details)
	}

	return
}
