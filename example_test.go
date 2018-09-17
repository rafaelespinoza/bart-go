package bart_test

import (
	"fmt"

	"github.com/rafaelespinoza/bart-go"
)

func ExampleAdvisoriesAPI() {
	c := new(bart.Client)
	var res interface{}
	var err error

	res, err = c.AdvisoriesAPI.RequestBSA()
	fmt.Println(res, err)

	res, err = c.AdvisoriesAPI.RequestTrainCount()
	fmt.Println(res, err)

	res, err = c.AdvisoriesAPI.RequestElevator()
	fmt.Println(res, err)
}

func ExampleRoutesAPI() {
	c := new(bart.Client)
	var res interface{}
	var err error

	// get route information based today's schedule
	res, err = c.RoutesAPI.RequestRoutesInfo("")
	fmt.Println(res, err)

	// get route information based on specific date's schedule
	res, err = c.RoutesAPI.RequestRoutesInfo("1/1/2018")
	fmt.Println(res, err)

	// get information on current routes
	res, err = c.RoutesAPI.RequestRoutes("")
	fmt.Println(res, err)
}

func ExampleStationsAPI() {
	c := new(bart.Client)
	var res interface{}
	var err error

	// request access information on MacArthur station
	res, err = c.StationsAPI.RequestStationAccess("mcar")
	fmt.Println(res, err)

	// be sure to pass in a valid station abbreviation. this example will return an error.
	res, err = c.StationsAPI.RequestStationAccess("nope")
	fmt.Println(res, err)

	// get detailed info on a specific station
	res, err = c.StationsAPI.RequestStationInfo("civc")
	fmt.Println(res, err)

	// get all the stations
	res, err = c.StationsAPI.RequestStations()
	fmt.Println(res, err)
}

func ExampleEstimatesAPI() {
	c := new(bart.Client)
	var res interface{}
	var err error

	// get real-time estimates for all stations
	res, err = c.EstimatesAPI.RequestETD("ALL", "", "")
	fmt.Println(res, err)

	// get real-time estimates for one station, all platforms, directions
	res, err = c.EstimatesAPI.RequestETD("cast", "", "")
	fmt.Println(res, err)

	// get real-time estimates for one station, specific platform, any direction
	res, err = c.EstimatesAPI.RequestETD("nbrk", "2", "")
	fmt.Println(res, err)

	// get real-time estimates for one station, any platform, specific direction
	res, err = c.EstimatesAPI.RequestETD("mcar", "", "S")
	fmt.Println(res, err)
}

func ExampleSchedulesAPI() {
	c := new(bart.Client)
	var res interface{}
	var err error

	// Use the quick planner to request departure from SF airport to Castro Valley station, at
	// current date, departing now with 2 trips, before & after.
	departure := bart.TripParams{"sfia", "cast", "", "", 2, 2, true}
	res, err = c.RequestDepartures(departure)
	fmt.Println(res, err)

	// Use the quick planner to request arrival from Embarcadero to Coliseum station, at
	// current date, arriving at 6:30pm with 3 trips, before & after.
	arrival := bart.TripParams{"embr", "cols", "", "6:30pm", 3, 3, true}
	res, err = c.SchedulesAPI.RequestArrivals(arrival)
	fmt.Println(res, err)

	res, err = c.RequestHolidaySchedules()
	fmt.Println(res, err)

	res, err = c.RequestAvailableSchedules()
	fmt.Println(res, err)

	res, err = c.RequestSpecialSchedules()
	fmt.Println(res, err)

	// Get station schedules for current date.
	res, err = c.RequestStationSchedules("mcar", "")
	fmt.Println(res, err)

	// Get station schedules for specific date.
	res, err = c.RequestStationSchedules("glen", "9/16/2018")
	fmt.Println(res, err)

	// Request schedule for route 12, on a Saturday
	res, err = c.RequestRouteSchedules(12, "sa", "", true)
	fmt.Println(res, err)
}

func ExampleTripParams() {
	p := bart.TripParams{"sfia", "cast", "2:08pm", "09/16/2018", 1, 2, true}
	fmt.Printf("%#v", p)
	// Output: bart.TripParams{Orig:"sfia", Dest:"cast", Time:"2:08pm", Date:"09/16/2018", Before:1, After:2, Legend:true}
}
