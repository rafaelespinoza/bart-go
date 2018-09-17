package bart

// RoutesAPI is a namespace for route information requests to routes at /route.aspx. See official
// docs at https://api.bart.gov/docs/route/.
type RoutesAPI struct{}

// RequestRoutesInfo requests detailed information for all routes. You probably want to request the current schedule on
// the current date, so pass in empty values for sched, date inputs. Details on the format of those inputs are in
// official docs. See official docs at https://api.bart.gov/docs/route/routeinfo.aspx.
func (a *RoutesAPI) RequestRoutesInfo(sched, date string) (res RoutesInfoResponse, err error) {
	params := map[string]string{"route": "all"}
	if sched != "" {
		params["sched"] = sched
	}
	if date != "" {
		params["date"] = date
	}

	err = requestAPI(
		"/route.aspx",
		"routeinfo",
		params,
		&res,
	)

	return
}

// RoutesInfoResponse is the shape of an API response.
type RoutesInfoResponse struct {
	Root struct {
		ResponseMetaData
		SchedNum int `json:"sched_num,string"`
		Data     struct {
			List []struct {
				Name        string
				Abbr        string
				RouteID     string
				Number      int `json:",string"`
				Origin      string
				Destination string
				Direction   string // 'North' or 'South'
				Hexcolor    string
				Color       string
				Holidays    int `json:",string"`
				NumStations int `json:"num_stations,string"`
				Config      struct {
					Stations []string `json:"station"`
				}
			} `json:"route"`
		} `json:"Routes"`
	}
}

// RequestRoutes requests (less) detailed information on current routes. If you only want current schedule on current
// date, just pass empty strings for the sched, date inputs. See official docs at https://api.bart.gov/docs/route/routes.aspx.
func (a *RoutesAPI) RequestRoutes(sched, date string) (res RoutesResponse, err error) {
	params := map[string]string{}
	if sched != "" {
		params["sched"] = sched
	}
	if date != "" {
		params["date"] = date
	}

	err = requestAPI(
		"/route.aspx",
		"routes",
		params,
		&res,
	)

	return
}

// RoutesResponse is the shape of an API response.
type RoutesResponse struct {
	Root struct {
		ResponseMetaData
		SchedNum int `json:"sched_num,string"`
		Data     struct {
			List []struct {
				Name     string
				Abbr     string
				RouteID  string
				Number   int `json:",string"`
				Hexcolor string
				Color    string
			} `json:"Route"`
		} `json:"Routes"`
	}
}
