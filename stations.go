package bart

type StationsAPI struct{}

// http://api.bart.gov/docs/stn/stnaccess.aspx
func (a *StationsAPI) RequestStationAccess(orig string) (res StationAccessResponse, err error) {
	params := map[string]string{"orig": orig}

	err = requestApi(
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

// http://api.bart.gov/docs/stn/stninfo.aspx
func (a *StationsAPI) RequestStationInfo(orig string) (res StationInfoResponse, err error) {
	params := map[string]string{"orig": orig}

	err = requestApi(
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

// http://api.bart.gov/docs/stn/stns.aspx
func (a *StationsAPI) RequestStations() (res StationsResponse, err error) {
	params := map[string]string{}

	err = requestApi(
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
