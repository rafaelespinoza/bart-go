package bart

import (
	"fmt"
	"strconv"
	"strings"
)

type EstimatesAPI struct{}

var (
	validPlatforms  = []string{"1", "2", "3", "4"}
	validDirections = []string{"N", "n", "S", "s"}
)

// https://api.bart.gov/docs/etd/etd.aspx
func (a *EstimatesAPI) RequestETD(orig, plat, dir string) (res EstimatesResponse, err error) {
	params := map[string]string{"orig": orig}
	allOrigins := strings.ToLower(orig) == "all"

	if _, err := validateStationAbbr(orig); !allOrigins && err != nil {
		return res, err
	}

	if !allOrigins && plat != "" {
		if p, e := validatePlatform(plat); e != nil {
			return res, e
		} else {
			params["plat"] = p
		}
	} else if !allOrigins && dir != "" {
		if d, e := validateDir(dir); e != nil {
			return res, e
		} else {
			params["dir"] = d
		}
	}

	err = requestApi(
		"/etd.aspx",
		"etd",
		params,
		&res,
	)

	return
}

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
					Minutes   estiMinute `json:",string"`
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

	if val, err := strconv.Atoi(str); err != nil {
		return err
	} else {
		*m = estiMinute(val)
		return nil
	}
}

func validatePlatform(plat string) (string, error) {
	if isPresent(plat, validPlatforms) {
		return plat, nil
	} else {
		err := fmt.Errorf("plat %q invalid. plat must be one of %v\n", plat, validPlatforms)
		return "", err
	}
}

func validateDir(dir string) (string, error) {
	if isPresent(dir, validDirections) {
		return dir, nil
	} else {
		err := fmt.Errorf("dir %q invalid. dir must be one of %v\n", dir, validDirections)
		return "", err
	}
}
