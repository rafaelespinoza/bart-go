package bart

import (
	"fmt"
	"strconv"
)

// SchedulesAPI is a namespace for schedule information requests to routes at
// /sched.aspx. See official docs at https://api.bart.gov/docs/sched/.
type SchedulesAPI struct{}

// RequestArrivals requests a trip plan based on arriving by the specified time.
// Inputs are specified in the TripParams type. See that type's documentation
// for details on requesting an arrival. See official docs at
// https://api.bart.gov/docs/sched/arrive.aspx.
func (a *SchedulesAPI) RequestArrivals(p TripParams) (res TripsResponse, err error) {
	params, err := p.validateMap()

	if err != nil {
		return res, err
	}

	err = requestAPI(
		"/sched.aspx",
		"arrive",
		params,
		&res,
	)

	return
}

// RequestDepartures requests a trip plan based on departing by the specified
// time. Inputs are specified in the TripParams type. See that type's
// documentation for details on requesting a departure. See official docs at
// https://api.bart.gov/docs/sched/depart.aspx.
func (a *SchedulesAPI) RequestDepartures(p TripParams) (res TripsResponse, err error) {
	params, err := p.validateMap()

	if err != nil {
		return res, err
	}

	err = requestAPI(
		"/sched.aspx",
		"depart",
		params,
		&res,
	)

	return
}

// TripsResponse is the shape of an API response.
type TripsResponse struct {
	Root struct {
		ResponseMetaData
		Message struct {
			CO2Emissions CDATASection `json:"co2_emissions"`
			Legend       string       `json:",omitempty"`
		}
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
						Order        int    `json:"@order,string"`
						TransferCode string `json:"@transfercode"`
						OrigDestTimeData
						Line             string  `json:"@line"`
						BikeFlag         boolish `json:"@bikeflag,string"`
						TrainHeadStation string  `json:"@trainHeadStation"`
						Load             int     `json:"@load,string"`
						TrainID          string  `json:"@trainId"`
						TrainIdx         int     `json:"@trainIdx,string"`
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
	params := make(map[string]string)

	err = requestAPI(
		"/sched.aspx",
		"holiday",
		params,
		&res,
	)

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
	params := make(map[string]string)

	err = requestAPI(
		"/sched.aspx",
		"scheds",
		params,
		&res,
	)

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
	params := make(map[string]string)

	err = requestAPI(
		"/sched.aspx",
		"special",
		params,
		&res,
	)

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

// RequestStationSchedules requests an entire daily schedule for the particular
// station specified. The orig param must be a 4-letter abbreviation for a
// station name. To request a schedule for a specific date, pass in a date
// formatted as "mm/dd/yyyy". Otherwise you can pass in "" to get today's
// schedule. See official docs at https://api.bart.gov/docs/sched/stnsched.aspx.
func (a *SchedulesAPI) RequestStationSchedules(orig, date string) (res StationSchedulesResponse, err error) {
	if _, err := validateStationAbbr(orig); err != nil {
		return res, err
	}

	params := map[string]string{"orig": orig}
	if date != "" {
		params["date"] = date
	}

	err = requestAPI(
		"/sched.aspx",
		"stnsched",
		params,
		&res,
	)

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
				Line             string  `json:"@line"`
				TrainHeadStation string  `json:"@trainHeadStation"`
				OrigTime         string  `json:"@origTime"`
				DestTime         string  `json:"@destTime"`
				TrainIdx         int     `json:"@trainIdx,string"`
				BikeFlag         boolish `json:"@bikeflag,string"`
				TrainID          string  `json:"@trainId"`
				Load             int     `json:"@load,string"`
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
func (a *SchedulesAPI) RequestRouteSchedules(
	route int,
	date string,
	time string,
	legend bool,
) (res RouteSchedulesResponse, err error) {
	params := map[string]string{"route": strconv.Itoa(route)}

	if date != "" {
		d, err := validateRouteSchedDate(date)
		if err != nil {
			return res, err
		}
		params["date"] = d
	}

	if time != "" {
		params["time"] = time
	}

	if legend {
		params["l"] = "1"
	}

	err = requestAPI(
		"/sched.aspx",
		"routesched",
		params,
		&res,
	)

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
					Station  string  `json:"@station"`
					Load     string  `json:"@load"`
					Level    string  `json:"@level"`
					BikeFlag boolish `json:"@bikeflag,string"`
				} `json:"stop"`
			} `json:"train"`
		} `json:"route"`
	}
}

// TripParams is a helper for two methods: RequestArrivals, RequestDepartures.
// It is used to manage inputs for those methods and to perform validation. The
// Orig and Dest fields are required and must be a 4-letter abbreviation for a
// station name. Passing in zero-values for both Before, After params is not
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

func (p TripParams) validateMap() (map[string]string, error) {
	params := map[string]string{}

	orig, err := validateStationAbbr(p.Orig)
	if err != nil {
		return params, err
	}
	params["orig"] = orig

	dest, err := validateStationAbbr(p.Dest)
	if err != nil {
		return params, err
	}
	params["dest"] = dest

	if p.Time != "" {
		params["time"] = p.Time
	}

	if p.Date != "" {
		params["date"] = p.Date
	}

	if p.Before == 0 && p.After == 0 {
		// API would return an empty string for value at `TripsResponse.Root.Data.Request`
		// in this case. I do not know how to handle that difference in type right now.
		return params, fmt.Errorf("before and after params cannot both == 0")
	}

	before, err := validateBeforeAfter(p.Before)
	if err != nil {
		return params, err
	}
	params["b"] = strconv.Itoa(before)

	after, err := validateBeforeAfter(p.After)
	if err != nil {
		return params, err
	}
	params["a"] = strconv.Itoa(after)

	if p.Legend {
		params["l"] = "1"
	}

	return params, nil
}

func validateBeforeAfter(val int) (int, error) {
	if val < 0 || val > 4 {
		err := fmt.Errorf("value %d invalid. param 'before' or 'after' must be >= 0 && <= 4", val)
		return 0, err
	}
	return val, nil
}

func validateRouteSchedDate(date string) (string, error) {
	// TODO: validate format `mm/dd/yyyy` or `wd|sa|su`
	return date, nil
}
