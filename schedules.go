package bart

import (
	"fmt"
	"strconv"
)

type SchedulesAPI struct{}

// https://api.bart.gov/docs/sched/arrive.aspx
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

// https://api.bart.gov/docs/sched/depart.aspx
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

type OrigDestTimeData struct {
	Origin       string `json:"@origin"`
	Destination  string `json:"@destination"`
	OrigTimeMin  string `json:"@origTimeMin"`
	OrigTimeDate string `json:"@origTimeDate"`
	DestTimeMin  string `json:"@destTimeMin"`
	DestTimeDate string `json:"@destTimeDate"`
}

// https://api.bart.gov/docs/sched/holiday.aspx
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

// https://api.bart.gov/docs/sched/scheds.aspx
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

// https://api.bart.gov/docs/sched/special.aspx
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

// https://api.bart.gov/docs/sched/stnsched.aspx
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

// https://api.bart.gov/docs/sched/routesched.aspx
func (a *SchedulesAPI) RequestRouteSchedules(
	route int,
	sched int,
	date string,
	time string,
	legend bool,
) (res RouteSchedulesResponse, err error) {
	params := map[string]string{"route": strconv.Itoa(route)}

	if sched != 0 {
		s, err := validateRouteSchedNum(sched)
		if err != nil {
			return res, err
		}
		params["sched"] = string(s)
	} else if date != "" {
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

func validateRouteSchedNum(sched int) (int, error) {
	// TODO: should be > 40, BART doesn't seem to have published schedules before 41
	return sched, nil
}

func validateRouteSchedDate(date string) (string, error) {
	// TODO: validate format `mm/dd/yyyy` or `wd|sa|su`
	return date, nil
}
