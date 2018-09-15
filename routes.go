package bart

// http://api.bart.gov/docs/route/routeinfo.aspx
func RequestRouteInfo(schedNum, date string) (res RouteInfoResponse, err error) {
	params := map[string]string{"route": "all"}
	if schedNum != "" {
		params["sched"] = schedNum
	}
	if date != "" {
		params["date"] = date
	}

	err = RequestApi(
		"/route.aspx",
		"routeinfo",
		params,
		&res,
	)

	return
}

type RouteInfoResponse struct {
	Root struct {
		ResponseMetaData
		SchedNum int                         `json:"sched_num,string"`
		Data     struct{ Route []RouteInfo } `json:"Routes"`
	}
}

type RouteInfo struct {
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
	Config      struct{ Station []string }
}

// http://api.bart.gov/docs/route/routes.aspx
func RequestRoutes(schedNum, date string) (res RoutesResponse, err error) {
	params := map[string]string{}
	if schedNum != "" {
		params["sched"] = schedNum
	}
	if date != "" {
		params["date"] = date
	}

	err = RequestApi(
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
		SchedNum int                     `json:"sched_num,string"`
		Data     struct{ Route []Route } `json:"Routes"`
	}
}

type Route struct {
	Name     string
	Abbr     string
	RouteId  string
	Number   int `json:",string"`
	Hexcolor string
	Color    string
}
