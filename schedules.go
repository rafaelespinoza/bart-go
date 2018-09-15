package bart

import (
	"fmt"
	"strconv"
)

// https://api.bart.gov/docs/sched/arrive.aspx
func RequestSchedInfoArrival(
	orig, dest string,
	time, date string,
	before, after int,
	legend bool,
) (res SchedInfoTripResponse, err error) {
	params, err := processScheduleParams(orig, dest, time, date, before, after, legend)

	if err != nil {
		return res, err
	}

	err = RequestApi(
		"/sched.aspx",
		"arrive",
		params,
		&res,
	)

	return
}

// https://api.bart.gov/docs/sched/depart.aspx
func RequestSchedInfoDeparture(
	orig, dest string,
	time, date string,
	before, after int,
	legend bool,
) (res SchedInfoTripResponse, err error) {
	params, err := processScheduleParams(orig, dest, time, date, before, after, legend)

	if err != nil {
		return res, err
	}

	err = RequestApi(
		"/sched.aspx",
		"depart",
		params,
		&res,
	)

	return
}

type SchedInfoTripResponse struct {
	Root struct {
		ResponseMetaData
		Message struct {
			Co2_emissions cDataSection
			Legend        string `json:",omitempty"`
		}
		Origin      string
		Destination string
		SchedNum    int      `json:"sched_num,string"`
		Data        Schedule `json:"schedule"`
	}
}

type Schedule struct {
	Date    string
	Time    string
	Before  int `json:",string"`
	After   int `json:",string"`
	Request struct {
		List []Trip `json:"Trip"`
	}
}

type Trip struct {
	OrigDestTimeData
	TripTime int       `json:"@tripTime,string"`
	Legs     []TripLeg `json:"leg"`
}

type TripLeg struct {
	Order        int    `json:"@order,string"`
	TransferCode string `json:"@transfercode"`
	OrigDestTimeData
	Line             string  `json:"@line"`
	BikeFlag         boolish `json:"@bikeflag,string"`
	TrainHeadStation string  `json:"@trainHeadStation"`
	Load             int     `json:"@load,string"`
	TrainId          string  `json:"@trainId"`
	TrainIdx         int     `json:"@trainIdx,string"`
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
func RequestSchedInfoHolidaySchedule() (res SchedInfoHolidayScheduleResponse, err error) {
	params := make(map[string]string)

	err = RequestApi(
		"/sched.aspx",
		"holiday",
		params,
		&res,
	)

	return
}

type SchedInfoHolidayScheduleResponse struct {
	Root struct {
		ResponseMetaData
		Data []struct {
			List []Holiday `json:"holiday"`
		} `json:"holidays"`
	}
}

type Holiday struct {
	Name         string
	Date         string
	ScheduleType string `json:"schedule_type"`
}

// https://api.bart.gov/docs/sched/scheds.aspx
func RequestSchedInfoSchedules() (res SchedInfoSchedulesResponse, err error) {
	params := make(map[string]string)

	err = RequestApi(
		"/sched.aspx",
		"scheds",
		params,
		&res,
	)

	return
}

type SchedInfoSchedulesResponse struct {
	Root struct {
		ResponseMetaData
		Data struct {
			List []ScheduleEdition `json:"schedule"`
		} `json:"schedules"`
	}
}

type ScheduleEdition struct {
	Id            int    `json:"@id,string"`
	EffectiveDate string `json:"@effectivedate"`
}

// https://api.bart.gov/docs/sched/special.aspx
func RequestSchedInfoSpecialSchedules() (res SchedInfoSpecialSchedulesResponse, err error) {
	params := make(map[string]string)

	err = RequestApi(
		"/sched.aspx",
		"special",
		params,
		&res,
	)

	return
}

type SchedInfoSpecialSchedulesResponse struct {
	Root struct {
		ResponseMetaData
		Data struct {
			List []SpecialSchedule `json:"special_schedule"`
		} `json:"special_schedules"`
	}
}

type SpecialSchedule struct {
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	Text           cDataSection
	Link           cDataSection
	Orig           string
	Dest           string
	DayOfWeek      string `json:"day_of_week"`
	RoutesAffected string `json:"routes_affected"`
}

// https://api.bart.gov/docs/sched/stnsched.aspx
func RequestSchedInfoStationSchedule(orig, date string) (res SchedInfoStationScheduleResponse, err error) {
	params := map[string]string{"orig": orig}
	if date != "" {
		params["date"] = date
	}

	err = RequestApi(
		"/sched.aspx",
		"stnsched",
		params,
		&res,
	)

	return
}

type SchedInfoStationScheduleResponse struct {
	Root struct {
		ResponseMetaData
		SchedNum int `json:"sched_num,string"`
		Data     struct {
			Name string
			Abbr string
			List []StationSchedule `json:"item"`
		} `json:"station"`
	}
}

type StationSchedule struct {
	Line             string  `json:"@line"`
	TrainHeadStation string  `json:"@trainHeadStation"`
	OrigTime         string  `json:"@origTime"`
	DestTime         string  `json:"@destTime"`
	TrainIdx         int     `json:"@trainIdx,string"`
	BikeFlag         boolish `json:"@bikeflag,string"`
	TrainId          string  `json:"@trainId"`
	Load             int     `json:"@load,string"`
}

// https://api.bart.gov/docs/sched/routesched.aspx
func RequestSchedInfoRouteSchedule(
	route, sched int,
	date, time string,
	legend bool,
) (res SchedInfoRouteScheduleResponse, err error) {
	params := map[string]string{"route": strconv.Itoa(route)}

	if sched != 0 {
		if s, e := validateRouteSchedNum(sched); e != nil {
			return res, e
		} else {
			params["sched"] = string(s)
		}
	} else if date != "" {
		if d, e := validateRouteSchedDate(date); e != nil {
			return res, e
		} else {
			params["date"] = d
		}
	}

	if time != "" {
		params["time"] = time
	}

	if legend {
		params["l"] = "1"
	}

	err = RequestApi(
		"/sched.aspx",
		"routesched",
		params,
		&res,
	)

	return
}

type SchedInfoRouteScheduleResponse struct {
	Root struct {
		ResponseMetaData
		SchedNum int `json:"sched_num,string"`
		Data     struct {
			List []TrainSchedule `json:"train"`
		} `json:"route"`
	}
}

type TrainSchedule struct {
	TrainId  string      `json:"@trainId"`
	TrainIdx int         `json:"@trainIdx,string"`
	Index    int         `json:"@index,string"`
	Stops    []TrainStop `json:"stop"`
}

type TrainStop struct {
	Station  string  `json:"@station"`
	Load     string  `json:"@load"`
	Level    string  `json:"@level"`
	BikeFlag boolish `json:"@bikeflag,string"`
}

func processScheduleParams(
	orig, dest string,
	time, date string,
	before, after int,
	legend bool,
) (map[string]string, error) {
	params := map[string]string{"orig": orig, "dest": dest}

	if time != "" {
		params["time"] = time
	}

	if date != "" {
		params["date"] = date
	}

	if before != 0 {
		if b, e := validateBeforeAfter(before); e != nil {
			return params, e
		} else {
			params["b"] = string(b)
		}
	}

	if after != 0 {
		if a, e := validateBeforeAfter(after); e != nil {
			return params, e
		} else {
			params["a"] = string(a)
		}
	}

	if legend {
		params["l"] = "1"
	}

	return params, nil
}

func validateBeforeAfter(val int) (int, error) {
	if val < 0 || val > 4 {
		msg := "value %d invalid. param 'before' or 'after' must be >= 0 && <= 4\n"
		err := fmt.Errorf(msg, val)
		return 0, err
	} else {
		return val, nil
	}
}

func validateRouteSchedNum(sched int) (int, error) {
	// TODO: should be > 40, BART doesn't seem to have published schedules before 41
	return sched, nil
}

func validateRouteSchedDate(date string) (string, error) {
	// TODO: validate format `mm/dd/yyyy` or `wd|sa|su`
	return date, nil
}
