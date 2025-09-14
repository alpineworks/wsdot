package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
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

	// Create a single unified vessels client that handles vessels, terminals, and fares
	vesselsClient, err := vessels.NewVesselsClient(wsdotClient)
	if err != nil {
		panic(err)
	}

	fmt.Println("=== Testing Unified Vessels Client (Terminals and Fares) ===")

	// 1. Test Cache Flush Date
	fmt.Println("\n1. Cache Flush Date:")
	cacheDate, err := vesselsClient.GetCacheFlushDate()
	if err != nil {
		fmt.Printf("Error getting cache flush date: %v\n", err)
	} else {
		fmt.Printf("Cache flushed at: %v\n", cacheDate)
		prettyPrint(cacheDate, *verbose, "Cache Flush Date")
	}

	// 2. Test Terminal Basics
	fmt.Println("\n2. Terminal Basics:")
	terminalBasics, err := vesselsClient.GetTerminalBasics()
	if err != nil {
		fmt.Printf("Error getting terminal basics: %v\n", err)
	} else {
		fmt.Printf("Found %d terminals\n", len(terminalBasics))
		if len(terminalBasics) > 0 {
			fmt.Printf("First terminal: %s (%s) - Reservable: %t\n",
				terminalBasics[0].TerminalName, terminalBasics[0].TerminalAbbrev,
				terminalBasics[0].Reservable)
		}
		prettyPrint(terminalBasics, *verbose, "Terminal Basics")
	}

	// 3. Test Terminal Wait Times
	fmt.Println("\n3. Terminal Wait Times:")
	waitTimes, err := vesselsClient.GetTerminalWaitTimes()
	if err != nil {
		fmt.Printf("Error getting terminal wait times: %v\n", err)
	} else {
		fmt.Printf("Found %d terminal wait times\n", len(waitTimes))
		for i, waitTime := range waitTimes {
			if i >= 3 { // Show only first 3 to avoid too much output
				fmt.Printf("... and %d more terminals\n", len(waitTimes)-3)
				break
			}
			fmt.Printf("- %s (%s): %d minutes\n",
				waitTime.TerminalName, waitTime.TerminalAbbrev, waitTime.WaitTime)
		}
		prettyPrint(waitTimes, *verbose, "Terminal Wait Times")
	}

	// 4. Test Terminal Locations
	fmt.Println("\n4. Terminal Locations:")
	locations, err := vesselsClient.GetTerminalLocations()
	if err != nil {
		fmt.Printf("Error getting terminal locations: %v\n", err)
	} else {
		fmt.Printf("Found %d terminal locations\n", len(locations))
		if len(locations) > 0 {
			fmt.Printf("Example: %s at %.4f°N, %.4f°W\n",
				locations[0].TerminalName, locations[0].Latitude, locations[0].Longitude)
		}
		prettyPrint(locations, *verbose, "Terminal Locations")
	}

	// 5. Test Fare Terminals (using tomorrow's date)
	fmt.Println("\n5. Fare Terminals:")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	fareTerminals, err := vesselsClient.GetTerminals(tomorrow)
	if err != nil {
		fmt.Printf("Error getting fare terminals: %v\n", err)
	} else {
		fmt.Printf("Found %d fare terminals for %s\n", len(fareTerminals), tomorrow)
		if len(fareTerminals) > 0 {
			fmt.Printf("Example: %s (%s) - Region: %d\n",
				fareTerminals[0].TerminalName, fareTerminals[0].TerminalAbbrev,
				fareTerminals[0].Region)
		}
		prettyPrint(fareTerminals, *verbose, "Fare Terminals")
	}

	// 7. Test Terminal Mates (if we have terminals)
	if len(fareTerminals) > 0 {
		fmt.Println("\n7. Terminal Mates:")
		firstTerminalID := fareTerminals[0].TerminalID
		mates, err := vesselsClient.GetTerminalMates(tomorrow, firstTerminalID)
		if err != nil {
			fmt.Printf("Error getting terminal mates for terminal %d: %v\n", firstTerminalID, err)
		} else {
			fmt.Printf("Terminal %s can reach %d destinations\n",
				fareTerminals[0].TerminalName, len(mates))
			for i, mate := range mates {
				if i >= 3 { // Show only first 3
					fmt.Printf("... and %d more destinations\n", len(mates)-3)
					break
				}
				fmt.Printf("- %s (%s)\n", mate.TerminalName, mate.TerminalAbbrev)
			}
			prettyPrint(mates, *verbose, "Terminal Mates")
		}

		// 8. Test Fare Line Items (if we have terminal mates)
		if len(mates) > 0 {
			fmt.Println("\n8. Fare Line Items:")
			departingID := strconv.Itoa(firstTerminalID)
			arrivingID := strconv.Itoa(mates[0].TerminalID)

			fareItems, err := vesselsClient.GetFareLineItemsBasic(tomorrow, departingID, arrivingID, false)
			if err != nil {
				fmt.Printf("Error getting fare line items: %v\n", err)
			} else {
				fmt.Printf("Found %d fare options from %s to %s\n",
					len(fareItems), fareTerminals[0].TerminalName, mates[0].TerminalName)
				for i, item := range fareItems {
					if i >= 5 { // Show only first 5
						fmt.Printf("... and %d more fare options\n", len(fareItems)-5)
						break
					}
					fmt.Printf("- %s: $%.2f\n", item.FareLineItemDescription, item.Cost)
				}
				prettyPrint(fareItems, *verbose, "Fare Line Items")

				// 9. Test Fare Total (for first fare item)
				if len(fareItems) > 0 {
					fmt.Println("\n9. Fare Total Calculation:")
					fareTotal, err := vesselsClient.GetFareTotal(tomorrow, departingID, arrivingID, false,
						fareItems[0].FareLineItemID, 1)
					if err != nil {
						fmt.Printf("Error getting fare total: %v\n", err)
					} else {
						fmt.Printf("Total cost for 1x %s: $%.2f\n",
							fareTotal.Description, fareTotal.TotalCost)
						prettyPrint(fareTotal, *verbose, "Fare Total")
					}
				}
			}
		}
	}

	// 10. Test specific terminal by ID (if we have terminals)
	if len(terminalBasics) > 0 {
		fmt.Println("\n10. Specific Terminal Details:")
		terminalID := terminalBasics[0].TerminalID

		// Get specific terminal basic info
		terminalDetail, err := vesselsClient.GetTerminalBasicByID(terminalID)
		if err != nil {
			fmt.Printf("Error getting terminal details for ID %d: %v\n", terminalID, err)
		} else {
			fmt.Printf("Terminal %s details: %s, %s %s\n",
				terminalDetail.TerminalName, terminalDetail.City,
				terminalDetail.State, terminalDetail.ZipCode)
		}

		// Get specific terminal wait time
		terminalWaitTime, err := vesselsClient.GetTerminalWaitTimeByID(terminalID)
		if err != nil {
			fmt.Printf("Error getting terminal wait time for ID %d: %v\n", terminalID, err)
		} else {
			fmt.Printf("Current wait time at %s: %d minutes\n",
				terminalWaitTime.TerminalName, terminalWaitTime.WaitTime)
		}

		// Get specific terminal location
		terminalLocation, err := vesselsClient.GetTerminalLocationByID(terminalID)
		if err != nil {
			fmt.Printf("Error getting terminal location for ID %d: %v\n", terminalID, err)
		} else {
			fmt.Printf("Terminal %s coordinates: %.4f°N, %.4f°W\n",
				terminalLocation.TerminalName, terminalLocation.Latitude, terminalLocation.Longitude)
		}
	}

	fmt.Println("\n=== All unified vessels client tests completed! ===")
}