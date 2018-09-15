package bart

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	validPlatforms  = []string{"1", "2", "3", "4"}
	validDirections = []string{"N", "n", "S", "s"}
)

// https://api.bart.gov/docs/etd/etd.aspx
func RequestEstimate(orig, plat, dir string) (res EstimateResponse, err error) {
	params := map[string]string{"orig": orig}
	allOrigins := strings.ToLower(orig) == "all"

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

	err = RequestApi(
		"/etd.aspx",
		"etd",
		params,
		&res,
	)

	return
}

type EstimateResponse struct {
	Root struct {
		ResponseMetaData
		Data []EstimateStation `json:"Station"`
	}
}

type EstimateStation struct {
	Name string
	Abbr string
	Etds []Etd `json:"etd"`
}

type Etd struct {
	Destination  string
	Abbreviation string
	Limited      string
	Estimates    []Estimate `json:"estimate"`
}

type Estimate struct {
	Minutes   estiMinute `json:",string"`
	Platform  int        `json:",string"`
	Direction string
	Length    int `json:",string"`
	Color     string
	Hexcolor  string
	BikeFlag  boolish `json:",string"`
	Delay     int     `json:",string"`
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
