package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"alpineworks.io/wsdot"
	"alpineworks.io/wsdot/traffic"
)

// prettyPrint outputs a struct as indented JSON if verbose is true
func prettyPrint(v interface{}, verbose bool, label string) {
	if verbose {
		fmt.Printf("\n--- %s (Full JSON) ---\n", label)
		jsonBytes, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonBytes))
		fmt.Println("--- End JSON ---")
	}
}

func main() {
	// Parse command line flags
	verbose := flag.Bool("verbose", false, "Print full JSON structures for all responses")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s                    # Run with summary output\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -verbose           # Run with full JSON output\n", os.Args[0])
	}
	flag.Parse()

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

	// Create a new Traffic client
	trafficClient, err := traffic.NewTrafficClient(wsdotClient)
	if err != nil {
		panic(err)
	}

	fmt.Println("=== Testing All Traffic API Endpoints ===")

	// 1. Test Border Crossings
	fmt.Println("\n1. Border Crossings:")
	borderCrossings, err := trafficClient.GetBorderCrossings()
	if err != nil {
		fmt.Printf("Error getting border crossings: %v\n", err)
	} else {
		fmt.Printf("Found %d border crossings\n", len(borderCrossings))
		if len(borderCrossings) > 0 {
			fmt.Printf("Example: %s - Wait time: %d minutes\n",
				borderCrossings[0].CrossingName, borderCrossings[0].WaitTime)
		}
		prettyPrint(borderCrossings, *verbose, "Border Crossings")
	}

	// 2. Test Bridge Clearances
	fmt.Println("\n2. Bridge Clearances:")
	bridgeClearances, err := trafficClient.GetBridgeClearances()
	if err != nil {
		fmt.Printf("Error getting bridge clearances: %v\n", err)
	} else {
		fmt.Printf("Found %d bridge clearances\n", len(bridgeClearances))
		if len(bridgeClearances) > 0 {
			fmt.Printf("Example: %s - Min clearance: %d inches\n",
				bridgeClearances[0].CrossingDescription,
				bridgeClearances[0].VerticalClearanceMinimumInches)
		}
		prettyPrint(bridgeClearances, *verbose, "Bridge Clearances")
	}

	// 3. Test Commercial Vehicle Restrictions
	fmt.Println("\n3. Commercial Vehicle Restrictions:")
	cvRestrictions, err := trafficClient.GetCommercialVehicleRestrictions()
	if err != nil {
		fmt.Printf("Error getting CV restrictions: %v\n", err)
	} else {
		fmt.Printf("Found %d commercial vehicle restrictions\n", len(cvRestrictions))
		if len(cvRestrictions) > 0 {
			fmt.Printf("Example: %s - %s\n",
				cvRestrictions[0].LocationName, cvRestrictions[0].RestrictionType)
		}
		prettyPrint(cvRestrictions, *verbose, "Commercial Vehicle Restrictions")
	}

	// 4. Test Highway Alerts
	fmt.Println("\n4. Highway Alerts:")
	alerts, err := trafficClient.GetHighwayAlerts()
	if err != nil {
		fmt.Printf("Error getting highway alerts: %v\n", err)
	} else {
		fmt.Printf("Found %d highway alerts\n", len(alerts))
		if len(alerts) > 0 {
			fmt.Printf("Example: %s - %s\n",
				alerts[0].HeadlineDescription, alerts[0].EventCategory)
		}
		prettyPrint(alerts, *verbose, "Highway Alerts")
	}

	// 5. Test Mountain Pass Conditions
	fmt.Println("\n5. Mountain Pass Conditions:")
	passConditions, err := trafficClient.GetMountainPassConditions()
	if err != nil {
		fmt.Printf("Error getting mountain pass conditions: %v\n", err)
	} else {
		fmt.Printf("Found %d mountain pass conditions\n", len(passConditions))
		if len(passConditions) > 0 {
			fmt.Printf("Example: %s - %s, %d°F\n",
				passConditions[0].MountainPassName,
				passConditions[0].WeatherCondition,
				passConditions[0].TemperatureInFahrenheit)
		}
		prettyPrint(passConditions, *verbose, "Mountain Pass Conditions")
	}

	// 6. Test Toll Rates
	fmt.Println("\n6. Toll Rates:")
	tollRates, err := trafficClient.GetTollRates()
	if err != nil {
		fmt.Printf("Error getting toll rates: %v\n", err)
	} else {
		fmt.Printf("Found %d toll rates\n", len(tollRates))
		if len(tollRates) > 0 {
			fmt.Printf("Example: %s - $%.2f\n",
				tollRates[0].TripName, tollRates[0].CurrentToll)
		}
		prettyPrint(tollRates, *verbose, "Toll Rates")
	}

	// 7. Test Traffic Flow
	fmt.Println("\n7. Traffic Flow:")
	trafficFlow, err := trafficClient.GetTrafficFlow()
	if err != nil {
		fmt.Printf("Error getting traffic flow: %v\n", err)
	} else {
		fmt.Printf("Found %d traffic flow readings\n", len(trafficFlow))
		if len(trafficFlow) > 0 {
			fmt.Printf("Example: %s - Flow: %v\n",
				trafficFlow[0].StationName, trafficFlow[0].FlowReadingValue)
		}
		prettyPrint(trafficFlow, *verbose, "Traffic Flow")
	}

	// 8. Test Travel Times
	fmt.Println("\n8. Travel Times:")
	travelTimes, err := trafficClient.GetTravelTimes()
	if err != nil {
		fmt.Printf("Error getting travel times: %v\n", err)
	} else {
		fmt.Printf("Found %d travel time routes\n", len(travelTimes))
		if len(travelTimes) > 0 {
			fmt.Printf("Example: %s - Current: %d min, Average: %d min\n",
				travelTimes[0].Name, travelTimes[0].CurrentTime, travelTimes[0].AverageTime)
		}
		prettyPrint(travelTimes, *verbose, "Travel Times")
	}

	// 9. Test Weather Information
	fmt.Println("\n9. Weather Information:")
	weatherInfo, err := trafficClient.GetCurrentWeatherInformation()
	if err != nil {
		fmt.Printf("Error getting weather information: %v\n", err)
	} else {
		fmt.Printf("Found %d weather stations\n", len(weatherInfo))
		if len(weatherInfo) > 0 {
			fmt.Printf("Example: %s - %.1f°F, Wind: %.1f mph %s\n",
				weatherInfo[0].StationName, weatherInfo[0].TemperatureInFahrenheit,
				weatherInfo[0].WindSpeedInMPH, weatherInfo[0].WindDirectionCardinal)
		}
		prettyPrint(weatherInfo, *verbose, "Weather Information")
	}

	// 10. Test Weather Stations
	fmt.Println("\n10. Weather Stations:")
	weatherStations, err := trafficClient.GetWeatherStations()
	if err != nil {
		fmt.Printf("Error getting weather stations: %v\n", err)
	} else {
		fmt.Printf("Found %d weather stations\n", len(weatherStations))
		if len(weatherStations) > 0 {
			fmt.Printf("Example: %s at %.4f°N, %.4f°W\n",
				weatherStations[0].StationName,
				weatherStations[0].Latitude, weatherStations[0].Longitude)
		}
		prettyPrint(weatherStations, *verbose, "Weather Stations")
	}

	fmt.Println("\n=== All Traffic endpoint tests completed! ===")
}