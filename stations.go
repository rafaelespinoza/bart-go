package bart

import (
	"fmt"
	"strings"
)

// StationsAPI is a namespace for stations information requests to routes at
// /stn.aspx. See official docs at https://api.bart.gov/docs/stn/.
type StationsAPI struct{}

// RequestStationAccess requests detailed information how to access the
// specified station as well as information about the neighborhood around the
// station. Pass in a 4-letter abbreviation for a station as the orig param. See
// official docs at https://api.bart.gov/docs/stn/stnaccess.aspx.
func (a *StationsAPI) RequestStationAccess(orig string) (res StationAccessResponse, err error) {
	if _, e := validateStationAbbr(orig); e != nil {
		return res, e
	}

	params := map[string]string{"orig": orig}

	err = requestAPI(
		"/stn.aspx",
		"stnaccess",
		params,
		&res,
	)

	return
}

// StationAccessResponse is the shape of an API response.
type StationAccessResponse struct {
	Root struct {
		ResponseMetaData
		Data struct {
			StationAccess struct {
				Name            string
				Abbr            string
				Entering        CDATASection
				Exiting         CDATASection
				FillTime        CDATASection
				CarShare        CDATASection
				Lockers         CDATASection
				BikeStationText CDATASection
				Destinations    CDATASection
				Link            string
				ParkingFlag     boolish `json:"@parking_flag,string"`
				BikeFlag        boolish `json:"@bike_flag,string"`
				BikeStation     boolish `json:"@bike_station_flag,string"`
				LockerFlag      boolish `json:"@locker_flag,string"`
			} `json:"Station"`
		} `json:"Stations"`
	}
}

// RequestStationInfo provides a detailed information about the specified
// station. Pass in a 4-letter abbreviation for a station as the orig param. See
// official docs at https://api.bart.gov/docs/stn/stninfo.aspx.
func (a *StationsAPI) RequestStationInfo(orig string) (res StationInfoResponse, err error) {
	if _, e := validateStationAbbr(orig); e != nil {
		return res, e
	}

	params := map[string]string{"orig": orig}

	err = requestAPI(
		"/stn.aspx",
		"stninfo",
		params,
		&res,
	)

	return
}

// StationInfoResponse is the shape of an API response.
type StationInfoResponse struct {
	Root struct {
		ResponseMetaData
		Data struct {
			StationInfo struct {
				Name           string
				Abbr           string
				Latitude       float32 `json:"gtfs_latitude,string"`
				Longitude      float32 `json:"gtfs_longitude,string"`
				Address        string
				City           string
				County         string
				State          string
				ZipCode        string
				NorthRoutes    struct{ Route []string }    `json:"north_routes"`
				SouthRoutes    struct{ Route []string }    `json:"south_routes"`
				NorthPlatforms struct{ Platform []string } `json:"north_platforms"`
				SouthPlatforms struct{ Platform []string } `json:"south_platforms"`
				PlatformInfo   string                      `json:"platform_info"`
				Intro          CDATASection
				CrossStreet    CDATASection `json:"cross_street"`
				Food           CDATASection
				Shopping       CDATASection
				Attraction     CDATASection
				Link           CDATASection
			} `json:"Station"`
		} `json:"Stations"`
	}
}

// RequestStations provides a list of all available stations. See official docs
// at https://api.bart.gov/docs/stn/stns.aspx.
func (a *StationsAPI) RequestStations() (res StationsResponse, err error) {
	params := map[string]string{}

	err = requestAPI(
		"/stn.aspx",
		"stns",
		params,
		&res,
	)

	return
}

// StationsResponse is the shape of an API response.
type StationsResponse struct {
	Root struct {
		ResponseMetaData
		Data struct {
			List []struct {
				Name      string
				Abbr      string
				Latitude  float32 `json:"gtfs_latitude,string"`
				Longitude float32 `json:"gtfs_longitude,string"`
				Address   string
				City      string
				County    string
				State     string
				ZipCode   string
			} `json:"station"`
		} `json:"stations"`
	}
}

var stationAbbrs = map[string]bool{
	"12TH": true,
	"16TH": true,
	"19TH": true,
	"24TH": true,
	"ANTC": true,
	"ASHB": true,
	"BALB": true,
	"BAYF": true,
	"CAST": true,
	"CIVC": true,
	"COLM": true,
	"COLS": true,
	"CONC": true,
	"DALY": true,
	"DBRK": true,
	"DELN": true,
	"DUBL": true,
	"EMBR": true,
	"FRMT": true,
	"FTVL": true,
	"GLEN": true,
	"HAYW": true,
	"LAFY": true,
	"LAKE": true,
	"MCAR": true,
	"MLBR": true,
	"MONT": true,
	"NBRK": true,
	"NCON": true,
	"OAKL": true,
	"ORIN": true,
	"PCTR": true,
	"PHIL": true,
	"PITT": true,
	"PLZA": true,
	"POWL": true,
	"RICH": true,
	"ROCK": true,
	"SANL": true,
	"SBRN": true,
	"SFIA": true,
	"SHAY": true,
	"SSAN": true,
	"UCTY": true,
	"WARM": true,
	"WCRK": true,
	"WDUB": true,
	"WOAK": true,
}

func validateStationAbbr(s string) (string, error) {
	u := strings.ToUpper(s)

	if _, ok := stationAbbrs[u]; ok {
		return s, nil
	}
	return "", fmt.Errorf("%q not a valid station abbreviation", s)
}
