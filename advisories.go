package bart

type AdvisoriesAPI struct{}

// http://api.bart.gov/docs/bsa/bsa.aspx
func (a *AdvisoriesAPI) RequestBSA() (res AdvisoriesBSAResponse, err error) {
	params := make(map[string]string)

	err = requestAPI(
		"/bsa.aspx",
		"bsa",
		params,
		&res,
	)

	return
}

type AdvisoriesBSAResponse struct {
	Root struct {
		ResponseMetaData
		Data []struct {
			Description CDATASection
			Type        string
			Posted      string
		}
	}
}

// http://api.bart.gov/docs/bsa/elev.aspx
func (a *AdvisoriesAPI) RequestElevator() (res AdvisoriesElevatorResponse, err error) {
	params := make(map[string]string)

	err = requestAPI(
		"/bsa.aspx",
		"elev",
		params,
		&res,
	)

	return
}

type AdvisoriesElevatorResponse struct {
	Root struct {
		ResponseMetaData
		Data []struct {
			Station     string
			Type        string
			Description CDATASection
			Posted      string
			Expires     string
		} `json:"bsa"`
	}
}

// http://api.bart.gov/docs/bsa/count.aspx
func (a *AdvisoriesAPI) RequestTrainCount() (res AdvisoriesTrainCountResponse, err error) {
	params := make(map[string]string)

	err = requestAPI(
		"/bsa.aspx",
		"count",
		params,
		&res,
	)

	return
}

type AdvisoriesTrainCountResponse struct {
	Root struct {
		ResponseMetaData
		Data int `json:"TrainCount,string"`
	}
}
