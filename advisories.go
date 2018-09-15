package bart

// http://api.bart.gov/docs/bsa/bsa.aspx
func RequestAdvisoryService() (res AdvisoryBsaResponse, err error) {
	params := make(map[string]string)

	err = RequestApi(
		"/bsa.aspx",
		"bsa",
		params,
		&res,
	)

	return
}

type AdvisoryBsaResponse struct {
	Root struct {
		ResponseMetaData
		Data []Bsa `json:"Bsa"`
	}
}

type Bsa struct {
	Description cDataSection
	Type        string
	Posted      string
}

// http://api.bart.gov/docs/bsa/elev.aspx
func RequestAdvisoryElevator() (res AdvisoryElevatorResponse, err error) {
	params := make(map[string]string)

	err = RequestApi(
		"/bsa.aspx",
		"elev",
		params,
		&res,
	)

	return
}

type AdvisoryElevatorResponse struct {
	Root struct {
		ResponseMetaData
		Data []ElevatorAdvisory `json:"Bsa"`
	}
}

type ElevatorAdvisory struct {
	Station     string
	Type        string
	Description cDataSection
	Posted      string
	Expires     string
}

// http://api.bart.gov/docs/bsa/count.aspx
func RequestAdvisoryTrainCount() (res AdvisoryTrainCountResponse, err error) {
	params := make(map[string]string)

	err = RequestApi(
		"/bsa.aspx",
		"count",
		params,
		&res,
	)

	return
}

type AdvisoryTrainCountResponse struct {
	Root struct {
		ResponseMetaData
		Data int `json:"TrainCount,string"`
	}
}
