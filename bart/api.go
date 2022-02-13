package bart

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

type apiRequest struct {
	route   string
	cmd     string
	options map[string][]string
}

func (p apiRequest) requestAPI(cc configuredClient, out interface{}) error {
	conf := cc.clientConf()

	values := make(url.Values)
	values.Set("cmd", p.cmd)
	values.Set("key", conf.Key)
	values.Set("json", "y")
	for key, vals := range p.options {
		for _, val := range vals {
			values.Add(key, val)
		}
	}

	uri := conf.baseURL + p.route + "?" + values.Encode()
	res, err := conf.HTTP.Get(uri)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// It seems like the BART API just started returning non-200 status codes
	// when there's an error, although this is not documented anywhere. Until
	// there is some documentation, assume there might be an error buried in the
	// response body.
	if err := checkAPIError(raw); err != nil {
		return err
	}

	return json.Unmarshal(raw, out)
}

func checkAPIError(in []byte) error {
	var body struct {
		Root struct {
			Message struct {
				Error interface{}
			}
		}
	}

	switch err := json.Unmarshal(in, &body).(type) {
	case nil:
		// Attempt to crack open the error message below.
		break
	case *json.UnmarshalTypeError:
		// Looks like most successful requests will have a Root.Message field
		// that's actually a string.
		if err.Field == "Root.Message" && err.Value == "string" {
			return nil
		}
		return err
	case *json.SyntaxError:
		// Handle any errors from the BART API that are formatted as XML. As of
		// 2018, the JSON API format is in beta and they never got around to
		// converting errors to JSON.
		//
		// As of 2021-08-15, it looks like they may have fixed this. But,
		// leaving this in here in just case they haven't fixed it everywhere.
		var xmlMessage struct {
			Text    string `xml:"message>error>text"`
			Details string `xml:"message>error>details"`
		}
		xmlParseErr := xml.Unmarshal(in, &xmlMessage)
		if xmlParseErr != nil {
			return xmlParseErr
		}
		return fmt.Errorf("error: %s. %s", xmlMessage.Text, xmlMessage.Details)
	default:
		return fmt.Errorf("error checking api error: %v", err)
	}

	switch val := body.Root.Message.Error.(type) {
	case nil:
		// Interpret this as a successful response.
		break
	case map[string]interface{}:
		// Most of the time, the error value is a struct with the fields: Text
		// and Details. However, accept other key-value pairs as well because
		// changes to the error shape seem to happen without announcement.
		// Hopefully, a decent error message can be built up this way.
		msgs := make([]string, 0)
		for k, v := range val {
			msgs = append(msgs, fmt.Sprintf("%s: %v", k, v))
		}
		return errors.New(strings.Join(msgs, ", "))
	case string:
		// Yep, sometimes it's a string.
		return errors.New(val)
	default:
		// Not expected, but just go with the flow.
		return errors.New(fmt.Sprintf("%s", val))
	}

	return nil
}
