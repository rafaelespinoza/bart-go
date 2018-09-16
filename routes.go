package bart

type RoutesAPI struct{}

// http://api.bart.gov/docs/route/routeinfo.aspx
func (a *RoutesAPI) RequestRoutesInfo(sched, date string) (res RoutesInfoResponse, err error) {
	params := map[string]string{"route": "all"}
	if sched != "" {
		params["sched"] = sched
	}
	if date != "" {
		params["date"] = date
	}

	err = requestApi(
		"/route.aspx",
		"routeinfo",
		params,
		&res,
	)

	return
}

type RoutesInfoResponse struct {
	Root struct {
		ResponseMetaData
		SchedNum int `json:"sched_num,string"`
		Data     struct {
			List []struct {
				Name        string
				Abbr        string
				RouteId     string
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

// http://api.bart.gov/docs/route/routes.aspx
func (a *RoutesAPI) RequestRoutes(sched, date string) (res RoutesResponse, err error) {
	params := map[string]string{}
	if sched != "" {
		params["sched"] = sched
	}
	if date != "" {
		params["date"] = date
	}

	err = requestApi(
		"/route.aspx",
		"routes",
		params,
		&res,
	)

	return
}

type RoutesResponse struct {
	Root struct {
		ResponseMetaData
		SchedNum int `json:"sched_num,string"`
		Data     struct {
			List []struct {
				Name     string
				Abbr     string
				RouteId  string
				Number   int `json:",string"`
				Hexcolor string
				Color    string
			} `json:"Route"`
		} `json:"Routes"`
	}
}
