# WSDOT Go Client

A Go client library for accessing Washington State Department of Transportation (WSDOT) APIs. This library provides convenient access to traffic, vessel, and camera data from WSDOT's public APIs.

## Installation

```bash
go get alpineworks.io/wsdot
```

## Quick Start

```go
package main

import (
    "fmt"
    "os"

    "alpineworks.io/wsdot"
    "alpineworks.io/wsdot/vessels"
    "alpineworks.io/wsdot/traffic"
)

func main() {
    // Create client with API key
    client, err := wsdot.NewWSDOTClient(
        wsdot.WithAPIKey(os.Getenv("WSDOT_API_KEY")),
    )
    if err != nil {
        panic(err)
    }

    // Use vessels API
    vesselsClient, _ := vessels.NewVesselsClient(client)
    stats, err := vesselsClient.GetVesselStats()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Found %d vessels\n", len(stats))

    // Use traffic API
    trafficClient, _ := traffic.NewTrafficClient(client)
    alerts, err := trafficClient.GetHighwayAlerts()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Found %d highway alerts\n", len(alerts))
}
```

## API Modules

### Vessels
Ferry system data including schedules, vessel information, and real-time locations:
- Vessel statistics and accommodations
- Terminal locations and information
- Route schedules and sailing times
- Historical vessel positions
- Wait times and fares

### Traffic
Road and highway information:
- Highway alerts and incidents
- Travel times between points
- Traffic Cameras
- Traffic flow data
- Weather stations and conditions
- Mountain pass conditions
- Bridge clearances
- Border crossing wait times
- Commercial vehicle restrictions
- Toll rates

## Configuration

The client requires a WSDOT API key, which can be obtained from the [WSDOT website](https://wsdot.wa.gov/traffic/api/).

```go
client, err := wsdot.NewWSDOTClient(
    wsdot.WithAPIKey("your-api-key"),
    wsdot.WithHTTPClient(customHTTPClient), // optional
)
```

## Examples

See the `example/` directory for complete working examples:
- `example/vessels/` - Ferry and vessel data examples
- `example/traffic/` - Traffic and road condition examples
- `example/fares/` - Ferry fare information examples

Run examples with:
```bash
export WSDOT_API_KEY="your-api-key"
go run example/vessels/main.go
go run example/traffic/main.go -verbose
```

## Requirements

- Go 1.22 or higher
- Valid WSDOT API key
