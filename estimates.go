package bart

import (
	"strconv"
)

// EstimatesAPI is a namespace for real-time information requests to /etd.aspx.
// See official docs at https://api.bart.gov/docs/etd/.
type EstimatesAPI struct{}

// RequestETD requests estimated departure time for specified station. The orig
// param must be a 4-letter abbreviation for a station name. Specify plat "1",
// "2", "3", "4" for a specific platform, or an empty string for all platforms.
// Specify dir "n" for north, "s" for south, or you can pass empty string to get
// both directions.  See official docs at
// https://api.bart.gov/docs/etd/etd.aspx.
func (a *EstimatesAPI) RequestETD(orig, plat, dir string) (res EstimatesResponse, err error) {
	params := &EstimateParams{
		Orig: orig,
		Plat: plat,
		Dir:  dir,
	}
	if err != nil {
		return
	}

	err = requestAPI(
		"/etd.aspx",
		"etd",
		params.toMap(),
		&res,
	)

	return
}

// An EstimateParams is a set of named parameters for requesting estimated
// departures in real time. Orig should be the 4-letter abbreviation for the
// name of the station. Plat should be "1", "2", "3", "4", or "all". Dir should
// be "n" for North, "s" for South, or an empty string for both directions.
type EstimateParams struct {
	Orig string
	Plat string
	Dir  string
}

func (p EstimateParams) toMap() map[string]string {
	params := map[string]string{"orig": p.Orig}
	if p.Dir != "" {
		params["dir"] = p.Dir
	}
	if p.Plat != "" {
		params["plat"] = p.Plat
	}
	return params
}

// RequestEstimate requests estimated departures for a station. It's just like
// the RequestETD method except it takes an EstimateParams value. See official
// docs at https://api.bart.gov/docs/etd/etd.aspx.
func (a *EstimatesAPI) RequestEstimate(p EstimateParams) (res EstimatesResponse, err error) {
	params := p.toMap()

	err = requestAPI(
		"/etd.aspx",
		"etd",
		params,
		&res,
	)

	return
}

// EstimatesResponse is the shape of an API response. One field, under the
// Estimates key is of the private type, estiMinute. It's there because
// zero-value is not "0", but "Leaving". To make it easier to deserialize, this
// package aliases "Leaving" to int 0.
type EstimatesResponse struct {
	Root struct {
		ResponseMetaData
		Data []struct {
			Name string
			Abbr string
			Etds []struct {
				Destination  string
				Abbreviation string
				Limited      string
				Estimates    []struct {
					Minutes   estiMinute `json:",string"` // effectively an int. Exists b/c "Leaving" == 0
					Platform  int        `json:",string"`
					Direction string
					Length    int `json:",string"`
					Color     string
					Hexcolor  string
					BikeFlag  boolish `json:",string"`
					Delay     int     `json:",string"`
				} `json:"estimate"`
			} `json:"etd"`
		} `json:"station"`
	}
}

type estiMinute int

func (m *estiMinute) UnmarshalJSON(data []byte) error {
	str := string(data)

	if str == "Leaving" {
		*m = estiMinute(0)
		return nil
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		return err
	}

	*m = estiMinute(val)
	return nil
}
