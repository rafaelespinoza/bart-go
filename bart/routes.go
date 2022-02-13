package bart

func initRoutesRequest(cmd, date string) (out apiRequest) {
	out.route = "/route.aspx"
	out.cmd = cmd
	out.options = make(map[string][]string)
	if date != "" {
		out.options["date"] = []string{date}
	}
	return
}

// RoutesAPI is a namespace for route information requests to routes at
// /route.aspx. See official docs at https://api.bart.gov/docs/route/.
type RoutesAPI struct {
	conf *Config
}

func (a *RoutesAPI) clientConf() *Config {
	if a != nil && a.conf != nil {
		return a.conf
	}
	return defaultClientConf
}

// RequestRoutesInfo requests detailed information for all routes. You probably
// want to request the current schedule on the current date, so pass in "" for
// date. Otherwise, format like "mm/dd/yyyy". See official docs at
// https://api.bart.gov/docs/route/routeinfo.aspx.
func (a *RoutesAPI) RequestRoutesInfo(date string) (res RoutesInfoResponse, err error) {
	params := initRoutesRequest("routeinfo", date)
	params.options["route"] = []string{"all"}

	err = params.requestAPI(a, &res)
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

// RequestRoutes requests (less) detailed information on current routes. If you
// only want current schedule on current date, just pass empty strings for date.
// See official docs at https://api.bart.gov/docs/route/routes.aspx.
func (a *RoutesAPI) RequestRoutes(date string) (res RoutesResponse, err error) {
	params := initRoutesRequest("routes", date)
	err = params.requestAPI(a, &res)
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
