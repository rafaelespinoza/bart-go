package bart

import "testing"

func TestTrips(t *testing.T) {
	runTest := func(t *testing.T, res TripsResponse, err error) {
		t.Helper()

		if err != nil {
			t.Error(err.Error())
			return
		}

		list := res.Root.Data.Request.List
		if len(list) < 1 {
			t.Error("no results at Root.Data.Request.List")
			return
		}

		for i := range list {
			if list[i].OrigTimeMin == "" {
				t.Error("expected non-empty value for OrigTimeMin")
			}
			if list[i].DestTimeMin == "" {
				t.Error("expected non-empty value for DestTimeMin")
			}
		}
	}

	t.Run("Arrivals", func(t *testing.T) {
		server := makeTestServer(t, stubHandler{
			expectedPath:     "/sched.aspx",
			expectedCmd:      "arrive",
			responseFilename: "testdata/schedules/trips_ok.json",
		})
		defer server.Close()

		var (
			params TripParams
			res    TripsResponse
			err    error
		)

		client := NewClient(nil)
		client.conf.baseURL = server.URL

		params = TripParams{Orig: "woak", Dest: "embr"}
		res, err = client.RequestArrivals(params)
		runTest(t, res, err)

		params = TripParams{Orig: "woak", Dest: "embr", Before: 0, After: 1}
		res, err = client.RequestArrivals(params)
		runTest(t, res, err)

		params = TripParams{Orig: "woak", Dest: "embr", Before: 0, After: 4}
		res, err = client.RequestArrivals(params)
		runTest(t, res, err)
	})

	t.Run("Departures", func(t *testing.T) {
		server := makeTestServer(t, stubHandler{
			expectedPath:     "/sched.aspx",
			expectedCmd:      "depart",
			responseFilename: "testdata/schedules/trips_ok.json",
		})
		defer server.Close()

		var (
			params TripParams
			res    TripsResponse
			err    error
		)

		client := NewClient(nil)
		client.conf.baseURL = server.URL

		params = TripParams{Orig: "woak", Dest: "embr"}
		res, err = client.RequestDepartures(params)
		runTest(t, res, err)

		params = TripParams{Orig: "woak", Dest: "embr", Before: 0, After: 1}
		res, err = client.RequestDepartures(params)
		runTest(t, res, err)

		params = TripParams{Orig: "woak", Dest: "embr", Before: 0, After: 4}
		res, err = client.RequestDepartures(params)
		runTest(t, res, err)
	})
}

func TestSpecialSchedules(t *testing.T) {
	t.Run("empty data", func(t *testing.T) {
		server := makeTestServer(t, stubHandler{
			expectedPath:     "/sched.aspx",
			expectedCmd:      "special",
			responseFilename: "testdata/schedules/special_schedules_empty.json",
		})
		defer server.Close()

		client := NewClient(nil)
		client.conf.baseURL = server.URL

		_, err := client.RequestSpecialSchedules()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		// should be able to parse w/o error
	})

	t.Run("non-empty data", func(t *testing.T) {
		server := makeTestServer(t, stubHandler{
			expectedPath:     "/sched.aspx",
			expectedCmd:      "special",
			responseFilename: "testdata/schedules/special_schedules_non_empty.json",
		})
		defer server.Close()

		client := NewClient(nil)
		client.conf.baseURL = server.URL

		out, err := client.RequestSpecialSchedules()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		val := out.Root.Data
		if len(val.List) < 1 {
			t.Fatal("Expected non-empty data")
		}
		if len(val.List[0].StartDate) < 1 {
			t.Errorf("Expected value for StartDate")
		}
		if len(val.List[0].EndDate) < 1 {
			t.Errorf("Expected value for EndDate")
		}
		if len(val.List[0].Text.Value) < 1 {
			t.Errorf("Expected value for Text.Value")
		}
		if len(val.List[0].RoutesAffected) < 1 {
			t.Errorf("Expected value for RoutesAffected")
		}
	})
}
