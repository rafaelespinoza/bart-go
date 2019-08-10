package bart

// AdvisoriesAPI is a namespace for advisory information requests to routes at
// /bsa.aspx. See official docs at https://api.bart.gov/docs/bsa/.
type AdvisoriesAPI struct {
	conf *Config
}

func (a *AdvisoriesAPI) clientConf() *Config {
	if a != nil && a.conf != nil {
		return a.conf
	}
	return defaultClientConf
}

// RequestBSA requests current advisory information. See official docs at
// https://api.bart.gov/docs/bsa/bsa.aspx.
func (a *AdvisoriesAPI) RequestBSA() (res AdvisoriesBSAResponse, err error) {
	err = requestAPI(
		a,
		"/bsa.aspx",
		"bsa",
		nil,
		&res,
	)

	return
}

// AdvisoriesBSAResponse is the shape of an API response.
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

// RequestElevator requests current elevator status information. See official
// docs at https://api.bart.gov/docs/bsa/elev.aspx.
func (a *AdvisoriesAPI) RequestElevator() (res AdvisoriesElevatorResponse, err error) {
	err = requestAPI(
		a,
		"/bsa.aspx",
		"elev",
		nil,
		&res,
	)

	return
}

// AdvisoriesElevatorResponse is the shape of an API response.
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

// RequestTrainCount requests the number of trains currently active in the
// system. See official docs at: https://api.bart.gov/docs/bsa/count.aspx.
func (a *AdvisoriesAPI) RequestTrainCount() (res AdvisoriesTrainCountResponse, err error) {
	err = requestAPI(
		a,
		"/bsa.aspx",
		"count",
		nil,
		&res,
	)

	return
}

// AdvisoriesTrainCountResponse is the shape of an API response.
type AdvisoriesTrainCountResponse struct {
	Root struct {
		ResponseMetaData
		Data int `json:"TrainCount,string"`
	}
}
