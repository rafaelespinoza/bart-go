package bart_test

import (
	"fmt"

	"github.com/rafaelespinoza/bart-go"
)

func ExampleAdvisoriesAPI() {
	c := new(bart.Client)
	var res interface{}
	var err error

	// Get current advisories.
	res, err = c.AdvisoriesAPI.RequestBSA()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Get current elevator status.
	res, err = c.AdvisoriesAPI.RequestElevator()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Get number of trains currently active.
	res, err = c.AdvisoriesAPI.RequestTrainCount()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Output:
	// bart.AdvisoriesBSAResponse
	// bart.AdvisoriesElevatorResponse
	// bart.AdvisoriesTrainCountResponse
}

func ExampleRoutesAPI() {
	c := new(bart.Client)
	var res interface{}
	var err error

	// Get information on current routes.
	res, err = c.RoutesAPI.RequestRoutes("")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Get route information based today's schedule.
	res, err = c.RoutesAPI.RequestRoutesInfo("")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Get route information based on specific date's schedule.
	res, err = c.RoutesAPI.RequestRoutesInfo("1/1/2018")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Output:
	// bart.RoutesResponse
	// bart.RoutesInfoResponse
	// bart.RoutesInfoResponse
}

func ExampleStationsAPI() {
	c := new(bart.Client)
	var res interface{}
	var err error

	// Request access information on MacArthur station.
	res, err = c.StationsAPI.RequestStationAccess("mcar")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Invalid station inputs will return an error.
	res, err = c.StationsAPI.RequestStationAccess("nope")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Get detailed info on a specific station.
	res, err = c.StationsAPI.RequestStationInfo("civc")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Get all the stations.
	res, err = c.StationsAPI.RequestStations()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Output:
	// bart.StationAccessResponse
	// error: Invalid orig. The orig station parameter NOPE is missing or invalid.
	// bart.StationAccessResponse
	// bart.StationInfoResponse
	// bart.StationsResponse
}

func ExampleEstimatesAPI() {
	c := new(bart.Client)
	var res interface{}
	var err error

	// Get real-time estimates for all stations.
	res, err = c.EstimatesAPI.RequestETD("ALL", "", "")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Get real-time estimates for a station, all platforms, directions.
	res, err = c.EstimatesAPI.RequestETD("cast", "", "")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Get real-time estimates for a station, specific platform, any direction.
	res, err = c.EstimatesAPI.RequestETD("nbrk", "2", "")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Get real-time estimates for a station, any platform, specific direction.
	res, err = c.EstimatesAPI.RequestETD("mcar", "", "S")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Output:
	// bart.EstimatesResponse
	// bart.EstimatesResponse
	// bart.EstimatesResponse
	// bart.EstimatesResponse
}

func ExampleSchedulesAPI() {
	c := new(bart.Client)
	var res interface{}
	var err error

	// Use the quick planner to request arrival from Embarcadero to Coliseum
	// station, at current date, arriving at 6:30pm with 3 trips, before &
	// after.
	res, err = c.SchedulesAPI.RequestArrivals(bart.TripParams{
		Orig:   "embr",
		Dest:   "cols",
		Time:   "6:30pm",
		Date:   "",
		Before: 3,
		After:  3,
		Legend: true,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Get currently available schedules.
	res, err = c.RequestAvailableSchedules()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Use the quick planner to request departure from SF airport to Castro
	// Valley station, at current date, departing now with 2 trips, before &
	// after.
	res, err = c.RequestDepartures(bart.TripParams{
		Orig:   "sfia",
		Dest:   "cast",
		Time:   "",
		Date:   "",
		Before: 2,
		After:  2,
		Legend: true,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Get info on upcoming BART holidays and what schedules run on those days.
	res, err = c.RequestHolidaySchedules()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Request schedule for route 12, on a Saturday.
	res, err = c.RequestRouteSchedules(12, "sa", "", true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Get special schedule notices in effect.
	res, err = c.RequestSpecialSchedules()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Get a station's schedules for the current date.
	res, err = c.RequestStationSchedules("mcar", "")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Get a station's schedules for specific date..
	res, err = c.RequestStationSchedules("glen", "8/14/2018")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T\n", res)

	// Output:
	// bart.TripsResponse
	// bart.AvailableSchedulesResponse
	// bart.TripsResponse
	// bart.HolidaySchedulesResponse
	// bart.RouteSchedulesResponse
	// bart.SpecialSchedulesResponse
	// bart.StationSchedulesResponse
	// bart.StationSchedulesResponse
}
