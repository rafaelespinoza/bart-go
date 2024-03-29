package bart

import (
	"encoding/json"
	"strconv"
)

func initSchedulesRequest(cmd string) (out apiRequest) {
	out.route = "/sched.aspx"
	out.cmd = cmd
	out.options = make(map[string][]string)
	return
}

// SchedulesAPI is a namespace for schedule information requests to routes at
// /sched.aspx. See official docs at https://api.bart.gov/docs/sched/.
type SchedulesAPI struct {
	conf *Config
}

func (a *SchedulesAPI) clientConf() *Config {
	if a != nil && a.conf != nil {
		return a.conf
	}
	return defaultClientConf
}

// RequestArrivals requests a trip plan based on arriving by the specified time.
// Inputs are specified in the TripParams type. See that type's documentation
// for details on requesting an arrival. See official docs at
// https://api.bart.gov/docs/sched/arrive.aspx.
func (a *SchedulesAPI) RequestArrivals(p TripParams) (res TripsResponse, err error) {
	params := p.initRequestParams("arrive")
	err = params.requestAPI(a, &res)
	return
}

// RequestDepartures requests a trip plan based on departing by the specified
// time. Inputs are specified in the TripParams type. See that type's
// documentation for details on requesting a departure. See official docs at
// https://api.bart.gov/docs/sched/depart.aspx.
func (a *SchedulesAPI) RequestDepartures(p TripParams) (res TripsResponse, err error) {
	params := p.initRequestParams("depart")
	err = params.requestAPI(a, &res)
	return
}

// TripsResponse is the shape of an API response.
type TripsResponse struct {
	Root struct {
		ResponseMetaData
		Message     interface{}
		Origin      string
		Destination string
		SchedNum    int `json:"sched_num,string"`
		Data        struct {
			Date    string
			Time    string
			Before  int `json:",string"`
			After   int `json:",string"`
			Request struct {
				List []struct {
					OrigDestTimeData
					TripTime int `json:"@tripTime,string"`
					Legs     []struct {
						OrigDestTimeData
						Order            int    `json:"@order,string"`
						Line             string `json:"@line"`
						BikeFlag         Bool   `json:"@bikeflag,string"`
						TrainHeadStation string `json:"@trainHeadStation"`
						Load             int    `json:"@load,string"`
					} `json:"leg"`
				} `json:"Trip"`
			}
		} `json:"schedule"`
	}
}

// OrigDestTimeData is an internal helper container, only meant to DRY up some
// type definitions.
type OrigDestTimeData struct {
	Origin       string `json:"@origin"`
	Destination  string `json:"@destination"`
	OrigTimeMin  string `json:"@origTimeMin"`
	OrigTimeDate string `json:"@origTimeDate"`
	DestTimeMin  string `json:"@destTimeMin"`
	DestTimeDate string `json:"@destTimeDate"`
}

// RequestHolidaySchedules requests information on the upcoming BART holidays,
// and what type of schedule will be run on those days.
// https://api.bart.gov/docs/sched/holiday.aspx.
func (a *SchedulesAPI) RequestHolidaySchedules() (res HolidaySchedulesResponse, err error) {
	params := initSchedulesRequest("holiday")
	err = params.requestAPI(a, &res)
	return
}

// HolidaySchedulesResponse is the shape of an API response.
type HolidaySchedulesResponse struct {
	Root struct {
		ResponseMetaData
		Data []struct {
			List []struct {
				Name         string
				Date         string
				ScheduleType string `json:"schedule_type"`
			} `json:"holiday"`
		} `json:"holidays"`
	}
}

// RequestAvailableSchedules requests information about the currently available
// schedules. See official docs at https://api.bart.gov/docs/sched/scheds.aspx.
func (a *SchedulesAPI) RequestAvailableSchedules() (res AvailableSchedulesResponse, err error) {
	params := initSchedulesRequest("scheds")
	err = params.requestAPI(a, &res)
	return
}

// AvailableSchedulesResponse is the shape of an API response.
type AvailableSchedulesResponse struct {
	Root struct {
		ResponseMetaData
		Data struct {
			List []struct {
				ID            int    `json:"@id,string"`
				EffectiveDate string `json:"@effectivedate"`
			} `json:"schedule"`
		} `json:"schedules"`
	}
}

// RequestSpecialSchedules requests information about all special schedule
// notices in effect. See official docs at
// https://api.bart.gov/docs/sched/special.aspx.
func (a *SchedulesAPI) RequestSpecialSchedules() (res SpecialSchedulesResponse, err error) {
	params := initSchedulesRequest("special")
	err = params.requestAPI(a, &res)
	return
}

// SpecialSchedulesResponse is the shape of an API response.
type SpecialSchedulesResponse struct {
	Root struct {
		ResponseMetaData
		Data struct {
			List []struct {
				StartDate      string `json:"start_date"`
				EndDate        string `json:"end_date"`
				StartTime      string `json:"start_time"`
				EndTime        string `json:"end_time"`
				Text           CDATASection
				Link           CDATASection
				Orig           string
				Dest           string
				DayOfWeek      string `json:"day_of_week"`
				RoutesAffected string `json:"routes_affected"`
			} `json:"special_schedule"`
		} `json:"special_schedules"`
	}
}

func (r *SpecialSchedulesResponse) UnmarshalJSON(in []byte) (err error) {
	type specialSchedulesResponseJSON SpecialSchedulesResponse
	var s specialSchedulesResponseJSON

	err = json.Unmarshal(in, &s)
	switch err.(type) {
	case nil:
		*r = SpecialSchedulesResponse{Root: s.Root}
		return
	case *json.UnmarshalTypeError:
		// This is *probably* the case where Root.Data = "". In order to
		// gracefully unmarshal the input, assume there's no useful data.
		var t SpecialSchedulesResponse
		t.Root.ResponseMetaData = s.Root.ResponseMetaData
		*r = t
		err = nil
		return
	default:
		return
	}
}

// RequestStationSchedules requests an entire daily schedule for the particular
// station specified. The orig param must be a 4-letter abbreviation for a
// station name. To request a schedule for a specific date, pass in a date
// formatted as "mm/dd/yyyy". Otherwise you can pass in "" to get today's
// schedule. See official docs at https://api.bart.gov/docs/sched/stnsched.aspx.
func (a *SchedulesAPI) RequestStationSchedules(orig, date string) (res StationSchedulesResponse, err error) {
	params := initSchedulesRequest("stnsched")
	params.options["orig"] = []string{orig}

	if date != "" {
		params.options["date"] = []string{date}
	}

	err = params.requestAPI(a, &res)
	return
}

// StationSchedulesResponse is the shape of an API response.
type StationSchedulesResponse struct {
	Root struct {
		ResponseMetaData
		SchedNum int `json:"sched_num,string"`
		Data     struct {
			Name string
			Abbr string
			List []struct {
				Line             string `json:"@line"`
				TrainHeadStation string `json:"@trainHeadStation"`
				OrigTime         string `json:"@origTime"`
				DestTime         string `json:"@destTime"`
				TrainIdx         int    `json:"@trainIdx,string"`
				BikeFlag         Bool   `json:"@bikeflag,string"`
				TrainID          string `json:"@trainId"`
				Load             int    `json:"@load,string"`
			} `json:"item"`
		} `json:"station"`
	}
}

// RequestRouteSchedules requests a full schedule for the specified route.
// Values for the route param must be one of 1-8, 11-12 or 19-20. Other inputs
// to this method default to the current values for current schedule today. To
// request specific details, such as the schedule on a certain day or another
// edition of the schedule pass in non-zero values as needed. See official docs
// at https://api.bart.gov/docs/sched/routesched.aspx.
func (a *SchedulesAPI) RequestRouteSchedules(route int, date string, time string, legend bool) (res RouteSchedulesResponse, err error) {
	params := initSchedulesRequest("routesched")
	params.options["route"] = []string{strconv.Itoa(route)}
	if date != "" {
		params.options["date"] = []string{date}
	}
	if time != "" {
		params.options["time"] = []string{time}
	}
	if legend {
		params.options["l"] = []string{"1"}
	}

	err = params.requestAPI(a, &res)
	return
}

// RouteSchedulesResponse is the shape of an API response.
type RouteSchedulesResponse struct {
	Root struct {
		ResponseMetaData
		SchedNum int `json:"sched_num,string"`
		Data     struct {
			List []struct {
				TrainID  string `json:"@trainId"`
				TrainIdx int    `json:"@trainIdx,string"`
				Index    int    `json:"@index,string"`
				Stops    []struct {
					Station  string `json:"@station"`
					Load     string `json:"@load"`
					Level    string `json:"@level"`
					BikeFlag Bool   `json:"@bikeflag,string"`
				} `json:"stop"`
			} `json:"train"`
		} `json:"route"`
	}
}

// TripParams is a helper for two methods: RequestArrivals, RequestDepartures.
// The Orig and Dest fields are required and must be a 4-letter abbreviation for
// a station name. Passing in zero-values for both Before, After params is not
// allowed, however you can pass a zero-value to one or the other. Details on
// the formatting of Time, Date params can be found in the official BART API
// docs. Most of the time you'd want to use the zero-value for Time, Data params
// anyways so you can fallback to the current time and current date.
type TripParams struct {
	Orig   string
	Dest   string
	Time   string
	Date   string
	Before int
	After  int
	Legend bool
}

func (p *TripParams) initRequestParams(cmd string) (out apiRequest) {
	out = initSchedulesRequest(cmd)
	out.options["orig"] = []string{p.Orig}
	out.options["dest"] = []string{p.Dest}

	if p.Time != "" {
		out.options["time"] = []string{p.Time}
	}

	if p.Date != "" {
		out.options["date"] = []string{p.Date}
	}

	if p.Legend {
		out.options["l"] = []string{"1"}
	}

	// Sending empty value for both `b`, `a` params (or `b=0`, `a=1`) returns an
	// object or empty string at TripsResponse.Root.Data.Request, which makes
	// unmarshaling a real pain. Avoid this case by omitting params.
	if (p.Before == 0 && p.After == 0) || (p.Before == 0 && p.After == 1) {
		return
	}
	// values for Before, After are fixed by BART API if they are outside of
	// acceptable range.
	out.options["b"] = []string{strconv.Itoa(p.Before)}
	out.options["a"] = []string{strconv.Itoa(p.After)}

	return
}
