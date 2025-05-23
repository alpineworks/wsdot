package main

import (
	"fmt"
	"os"

	"alpineworks.io/wsdot"
	"alpineworks.io/wsdot/ferries"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("API_KEY environment variable is required")
	}

	// Create a new WSDOT client
	wsdotClient, err := wsdot.NewWSDOTClient(
		wsdot.WithAPIKey(apiKey),
	)

	if err != nil {
		panic(err)
	}

	// Create a new Ferries client
	ferriesClient, err := ferries.NewFerriesClient(wsdotClient)
	if err != nil {
		panic(err)
	}

	// Get the vessel basics
	vessels, err := ferriesClient.GetVesselBasics()
	if err != nil {
		panic(err)
	}

	if len(vessels) > 0 {
		fmt.Println(vessels[0].VesselName)
		fmt.Println(vessels[0].VesselID)
		fmt.Println(vessels[0].Class.ClassName)
	}

	// Get the vessel locations
	vesselLocations, err := ferriesClient.GetVesselLocations()
	if err != nil {
		panic(err)
	}

	if len(vesselLocations) > 0 {
		fmt.Println(vesselLocations[1].VesselName)
		fmt.Printf("%f°N, %f°W\n", vesselLocations[1].Latitude, vesselLocations[1].Longitude)
	}

	// get the route schedule
	schedules, err := ferriesClient.GetRouteSchedules()
	if err != nil {
		panic(err)
	}

	for _, schedule := range schedules {
		fmt.Printf("Schedule ID: %d, Route ID: %d, Description: %s\n", schedule.ScheduleID, schedule.RouteID, schedule.Description)
		allSailings, err := ferriesClient.GetSchedulesTodayByRouteID(schedule.RouteID, false)
		if err != nil {
			panic(err)
		}
		fmt.Printf("last sailing: %s\n", allSailings.TerminalCombos[0].Times[len(allSailings.TerminalCombos[0].Times)-1].DepartingTime)
	}

}
