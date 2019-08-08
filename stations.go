package bart

import "net/url"

// StationsAPI is a namespace for stations information requests to routes at
// /stn.aspx. See official docs at https://api.bart.gov/docs/stn/.
type StationsAPI struct {
	conf *Config
}

func (a *StationsAPI) clientConf() *Config {
	if a != nil && a.conf != nil {
		return a.conf
	}
	return defaultClientConf
}

// RequestStationAccess requests detailed information how to access the
// specified station as well as information about the neighborhood around the
// station. Pass in a 4-letter abbreviation for a station as the orig param. See
// official docs at https://api.bart.gov/docs/stn/stnaccess.aspx.
func (a *StationsAPI) RequestStationAccess(orig string) (res StationAccessResponse, err error) {
	params := url.Values{}
	params.Set("orig", orig)

	err = requestAPI(
		a,
		"/stn.aspx",
		"stnaccess",
		&params,
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
	params := url.Values{}
	params.Set("orig", orig)

	err = requestAPI(
		a,
		"/stn.aspx",
		"stninfo",
		&params,
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
	err = requestAPI(
		a,
		"/stn.aspx",
		"stns",
		nil,
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
