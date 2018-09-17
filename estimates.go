package bart

import (
	"fmt"
	"strconv"
	"strings"
)

// EstimatesAPI is a namespace for real-time information requests to /etd.aspx. See official docs at
// https://api.bart.gov/docs/etd/.
type EstimatesAPI struct{}

var (
	validPlatforms  = []string{"1", "2", "3", "4"}
	validDirections = []string{"N", "n", "S", "s"}
)

// RequestETD requests estimated departure time for specified station. The orig param must be a 4-letter abbreviation
// for a station name. Specify plat "1", "2", "3", "4" for a specific platform, or an empty string for all platforms.
// Specify dir "n" for north, "s" for south, or you can pass empty string to get both directions.  See official docs at
// https://api.bart.gov/docs/etd/etd.aspx.
func (a *EstimatesAPI) RequestETD(orig, plat, dir string) (res EstimatesResponse, err error) {
	params := map[string]string{"orig": orig}
	allOrigins := strings.ToLower(orig) == "all"

	if _, err := validateStationAbbr(orig); !allOrigins && err != nil {
		return res, err
	}

	if !allOrigins && plat != "" {
		p, e := validatePlatform(plat)
		if e != nil {
			return res, e
		}
		params["plat"] = p
	} else if !allOrigins && dir != "" {
		d, e := validateDir(dir)
		if e != nil {
			return res, e
		}
		params["dir"] = d
	}

	err = requestAPI(
		"/etd.aspx",
		"etd",
		params,
		&res,
	)

	return
}

// EstimatesResponse is the shape of an API response. One field, under the Estimates key is of the private type,
// estiMinute. It's there because zero-value is not "0", but "Leaving". To make it easier to deserialize, this package
// aliases "Leaving" to int 0.
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

func validatePlatform(plat string) (string, error) {
	if isPresent(plat, validPlatforms) {
		return plat, nil
	}

	err := fmt.Errorf("plat %q invalid, plat must be one of %v", plat, validPlatforms)
	return "", err
}

func validateDir(dir string) (string, error) {
	if isPresent(dir, validDirections) {
		return dir, nil
	}

	err := fmt.Errorf("dir %q invalid. dir must be one of %v", dir, validDirections)
	return "", err
}
