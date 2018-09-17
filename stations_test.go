package bart

import (
	"strings"
	"testing"
)

func TestValidateStationAbbr(t *testing.T) {
	type testCase struct {
		input  string
		expStr string
		expErr bool // do we expect an error?
	}

	validCases := make([]testCase, len(stationAbbrs)*2)
	i := 0
	for name := range stationAbbrs {
		j := i * 2
		lo := strings.ToLower(name)
		hi := strings.ToUpper(name)
		validCases[j] = testCase{lo, lo, false}
		validCases[j+1] = testCase{hi, hi, false}
		i++
	}

	tables := append(
		validCases,
		testCase{
			"BOOF",
			"",
			true,
		},
	)

	for _, test := range tables {
		str, err := validateStationAbbr(test.input)

		if str != test.expStr {
			t.Errorf("actual %s != expected %s", str, test.expStr)
		}

		errorPresent := err != nil
		if errorPresent != test.expErr {
			t.Errorf("error actually present?: %t != expected error presence %t", err, test.expErr)
		}
	}
}
