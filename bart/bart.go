// Package bart is a wrapper for the BART API. It works with JSON, which is
// still in beta at the time of this writing. See the official BART docs for
// information https://api.bart.gov/docs/overview/index.aspx
//
// The API request example functions are written with vague regard to the
// output. However, to execute the examples with go test, we must specify some
// kind of output. Making real requests to the BART API in tests (as opposed to
// stubbing responses) is meant to expose unexpected errors or panics, which
// should help make response handling better in the long run.
package bart

import "strconv"

// Bool is like a regular bool, but unmarshals the "boolish-looking" JSON values
// received from the BART API using strconv.ParseBool.
type Bool bool

func (b *Bool) UnmarshalJSON(data []byte) error {
	val, err := strconv.ParseBool(string(data))
	if err != nil {
		return err
	}
	*b = Bool(val)
	return nil
}

// Minute is like a minute, but on BART. It actually exists because the
// Real-Time Estimates API says "Leaving" when there are 0 minutes remaining
// until departure.
type Minute int

func (m *Minute) UnmarshalJSON(in []byte) error {
	data := string(in)
	if data == "Leaving" {
		*m = Minute(0)
		return nil
	}

	val, err := strconv.Atoi(data)
	if err != nil {
		return err
	}

	*m = Minute(val)
	return nil
}
