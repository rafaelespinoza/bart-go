package bart

// http://api.bart.gov/docs/stn/stnaccess.aspx
func RequestStationAccess(orig string) (res StationAccessResponse, err error) {
	params := map[string]string{"orig": orig}

	err = RequestApi(
		"/stn.aspx",
		"stnaccess",
		params,
		&res,
	)

	return
}

type StationAccessResponse struct {
	Root struct {
		ResponseMetaData
		Data struct {
			StationAccess `json:"Station"`
		} `json:"Stations"`
	}
}

type StationAccess struct {
	Name            string
	Abbr            string
	Entering        cDataSection
	Exiting         cDataSection
	FillTime        cDataSection
	CarShare        cDataSection
	Lockers         cDataSection
	BikeStationText cDataSection
	Destinations    cDataSection
	Link            string
	ParkingFlag     boolish `json:"@parking_flag,string"`
	BikeFlag        boolish `json:"@bike_flag,string"`
	BikeStation     boolish `json:"@bike_station_flag,string"`
	LockerFlag      boolish `json:"@locker_flag,string"`
}

// http://api.bart.gov/docs/stn/stninfo.aspx
func RequestStationInfo(orig string) (res StationInfoResponse, err error) {
	params := map[string]string{"orig": orig}

	err = RequestApi(
		"/stn.aspx",
		"stninfo",
		params,
		&res,
	)

	return
}

type StationInfoResponse struct {
	Root struct {
		ResponseMetaData
		Data struct {
			StationInfo `json:"Station"`
		} `json:"Stations"`
	}
}

type StationInfo struct {
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
	NorthPlatforms struct{ Platform []string } `json:"north_platforms"` // TODO: interpret as []int
	SouthPlatforms struct{ Platform []string } `json:"south_platforms"` // TODO: interpret as []int
	PlatformInfo   string                      `json:"platform_info"`
	Intro          cDataSection
	CrossStreet    cDataSection `json:"cross_street"`
	Food           cDataSection
	Shopping       cDataSection
	Attraction     cDataSection
	Link           cDataSection
}

// http://api.bart.gov/docs/stn/stns.aspx
func RequestStations() (res StationsResponse, err error) {
	params := map[string]string{}

	err = RequestApi(
		"/stn.aspx",
		"stns",
		params,
		&res,
	)

	return
}

type StationsResponse struct {
	Root struct {
		ResponseMetaData
		Data struct{ Station []Station } `json:"Stations"`
	}
}

type Station struct {
	Name      string
	Abbr      string
	Latitude  float32 `json:"gtfs_latitude,string"`
	Longitude float32 `json:"gtfs_longitude,string"`
	Address   string
	City      string
	County    string
	State     string
	ZipCode   string
}
