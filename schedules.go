package bart

import (
	"fmt"
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
