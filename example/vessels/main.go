package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"alpineworks.io/wsdot"
	"alpineworks.io/wsdot/vessels"
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

	// Create a new Vessels client
	vesselsClient, err := vessels.NewVesselsClient(wsdotClient)
	if err != nil {
		panic(err)
	}

	fmt.Println("=== Testing All Vessels API Endpoints ===")

	// 1. Test Cache Flush Date
	fmt.Println("\n1. Cache Flush Date:")
	cacheDate, err := vesselsClient.GetCacheFlushDate()
	if err != nil {
		fmt.Printf("Error getting cache flush date: %v\n", err)
	} else {
		fmt.Printf("Cache flushed at: %v\n", cacheDate)
		prettyPrint(cacheDate, *verbose, "Cache Flush Date")
	}

	// 2. Test Vessel Basics (all vessels)
	fmt.Println("\n2. Vessel Basics (all):")
	vesselBasics, err := vesselsClient.GetVesselBasics()
	if err != nil {
		fmt.Printf("Error getting vessel basics: %v\n", err)
	} else if len(vesselBasics) > 0 {
		fmt.Printf("Found %d vessels\n", len(vesselBasics))
		fmt.Printf("First vessel: %s (ID: %d, Class: %s)\n", 
			vesselBasics[0].VesselName, vesselBasics[0].VesselID, vesselBasics[0].Class.ClassName)
		prettyPrint(vesselBasics, *verbose, "All Vessel Basics")
		
		// 3. Test Vessel Basic by ID (using first vessel)
		fmt.Println("\n3. Vessel Basic by ID:")
		vesselBasic, err := vesselsClient.GetVesselBasicByID(vesselBasics[0].VesselID)
		if err != nil {
			fmt.Printf("Error getting vessel basic by ID: %v\n", err)
		} else {
			fmt.Printf("Vessel: %s (Owned by WSF: %t)\n", vesselBasic.VesselName, vesselBasic.OwnedByWSF)
			prettyPrint(vesselBasic, *verbose, "Single Vessel Basic")
		}
	}

	// 4. Test Vessel Locations (all vessels)
	fmt.Println("\n4. Vessel Locations (all):")
	vesselLocations, err := vesselsClient.GetVesselLocations()
	if err != nil {
		fmt.Printf("Error getting vessel locations: %v\n", err)
	} else if len(vesselLocations) > 0 {
		fmt.Printf("Found %d vessel locations\n", len(vesselLocations))
		for i, loc := range vesselLocations {
			if i >= 3 { // Show only first 3 to avoid too much output
				fmt.Printf("... and %d more vessels\n", len(vesselLocations)-3)
				break
			}
			fmt.Printf("- %s: %.4f°N, %.4f°W (Speed: %.1f knots)\n", 
				loc.VesselName, loc.Latitude, loc.Longitude, loc.Speed)
		}
		
		// 5. Test Vessel Location by ID (using first vessel)
		if len(vesselLocations) > 0 {
			fmt.Println("\n5. Vessel Location by ID:")
			vesselLocation, err := vesselsClient.GetVesselLocationByID(vesselLocations[0].VesselID)
			if err != nil {
				fmt.Printf("Error getting vessel location by ID: %v\n", err)
			} else {
				fmt.Printf("%s at dock: %t, In service: %t\n", 
					vesselLocation.VesselName, vesselLocation.AtDock, vesselLocation.InService)
			}
		}
	}

	// 6. Test Vessel Accommodations (all vessels)
	fmt.Println("\n6. Vessel Accommodations (all):")
	accommodations, err := vesselsClient.GetVesselAccommodations()
	if err != nil {
		fmt.Printf("Error getting vessel accommodations: %v\n", err)
	} else if len(accommodations) > 0 {
		fmt.Printf("Found %d vessel accommodations\n", len(accommodations))
		fmt.Printf("Example - %s: Car deck restroom: %t, Public WiFi: %t, ADA: %t\n",
			accommodations[0].VesselName, accommodations[0].CarDeckRestroom, 
			accommodations[0].PublicWifi, accommodations[0].ADAAccessible)
		prettyPrint(accommodations, *verbose, "All Vessel Accommodations")
		
		// 7. Test Vessel Accommodation by ID
		fmt.Println("\n7. Vessel Accommodation by ID:")
		accommodation, err := vesselsClient.GetVesselAccommodationByID(accommodations[0].VesselID)
		if err != nil {
			fmt.Printf("Error getting accommodation by ID: %v\n", err)
		} else {
			fmt.Printf("%s: Main cabin galley: %t, Car deck shelter: %t, Elevator: %t\n",
				accommodation.VesselName, accommodation.MainCabinGalley, 
				accommodation.CarDeckShelter, accommodation.Elevator)
		}
	}

	// 8. Test Vessel Stats (all vessels)
	fmt.Println("\n8. Vessel Stats (all):")
	stats, err := vesselsClient.GetVesselStats()
	if err != nil {
		fmt.Printf("Error getting vessel stats: %v\n", err)
	} else if len(stats) > 0 {
		fmt.Printf("Found %d vessel stats\n", len(stats))
		fmt.Printf("Example - %s: Built %d, %d HP, %d knots, Fast Ferry: %t\n",
			stats[0].VesselName, stats[0].YearBuilt, stats[0].Horsepower, 
			stats[0].SpeedInKnots, stats[0].FastFerry)
		
		// 9. Test Vessel Stats by ID
		fmt.Println("\n9. Vessel Stats by ID:")
		stat, err := vesselsClient.GetVesselStatsByID(stats[0].VesselID)
		if err != nil {
			fmt.Printf("Error getting stats by ID: %v\n", err)
		} else {
			fmt.Printf("%s: Length: %s, Beam: %s, Tonnage: %d, Passenger Only: %t\n",
				stat.VesselName, stat.Length, stat.Beam, stat.Tonnage, stat.PassengerOnly)
			prettyPrint(stat, *verbose, "Single Vessel Stats (Most Detailed)")
		}
	}

	// 10. Test Vessel History (all vessels)
	fmt.Println("\n10. Vessel History (all):")
	history, err := vesselsClient.GetVesselHistory()
	if err != nil {
		fmt.Printf("Error getting vessel history: %v\n", err)
	} else {
		fmt.Printf("Found %d historical records\n", len(history))
		if len(history) > 0 {
			fmt.Printf("Latest record: %s at %.4f°N, %.4f°W\n",
				history[0].VesselName, history[0].Latitude, history[0].Longitude)
		}
	}

	// 11. Test Vessel History by Name and Date Range
	fmt.Println("\n11. Vessel History by Name and Date Range:")
	if len(vesselBasics) > 0 {
		// Get history for the first vessel from yesterday to today
		startDate := time.Now().AddDate(0, 0, -1) // Yesterday
		endDate := time.Now()                     // Today
		vesselName := vesselBasics[0].VesselName
		
		historyRange, err := vesselsClient.GetVesselHistoryByNameAndDateRange(vesselName, startDate, endDate)
		if err != nil {
			fmt.Printf("Error getting vessel history by date range: %v\n", err)
		} else {
			fmt.Printf("Found %d records for %s from %s to %s\n", 
				len(historyRange), vesselName, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
		}
	}

	// 12. Test Route Schedules (Schedule API)
	fmt.Println("\n12. Route Schedules:")
	schedules, err := vesselsClient.GetRouteSchedules()
	if err != nil {
		fmt.Printf("Error getting route schedules: %v\n", err)
	} else if len(schedules) > 0 {
		fmt.Printf("Found %d route schedules\n", len(schedules))
		fmt.Printf("Example - Route %d: %s\n", schedules[0].RouteID, schedules[0].Description)
		
		// 13. Test Schedule Today by Route ID
		fmt.Println("\n13. Today's Schedule for Route:")
		allSailings, err := vesselsClient.GetSchedulesTodayByRouteID(schedules[0].RouteID, false)
		if err != nil {
			fmt.Printf("Error getting today's schedule: %v\n", err)
		} else {
			fmt.Printf("Schedule: %s, Season: %d\n", allSailings.ScheduleName, allSailings.ScheduleSeason)
			if len(allSailings.TerminalCombos) > 0 && len(allSailings.TerminalCombos[0].Times) > 0 {
				lastSailing := allSailings.TerminalCombos[0].Times[len(allSailings.TerminalCombos[0].Times)-1]
				fmt.Printf("Last sailing departing time: %v\n", lastSailing.DepartingTime)
			}
		}
	}

	fmt.Println("\n=== All endpoint tests completed! ===")
}
