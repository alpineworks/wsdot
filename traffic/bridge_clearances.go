package traffic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"alpineworks.io/wsdot"
)

const (
	getBridgeClearancesAsJsonURL = "http://wsdot.wa.gov/Traffic/api/Bridges/ClearanceREST.svc/GetClearancesAsJson"
)

// BridgeClearance represents a bridge clearance measurement from WSDOT
type BridgeClearance struct {
	APILastUpdate                      string  `json:"APILastUpdate"`
	BridgeNumber                       string  `json:"BridgeNumber"`
	ControlEntityGuid                  string  `json:"ControlEntityGuid"`
	CrossingDescription                string  `json:"CrossingDescription"`
	CrossingLocationId                 int     `json:"CrossingLocationId"`
	CrossingRecordGuid                 string  `json:"CrossingRecordGuid"`
	InventoryDirection                 string  `json:"InventoryDirection"`
	Latitude                           float64 `json:"Latitude"`
	LocationGuid                       string  `json:"LocationGuid"`
	Longitude                          float64 `json:"Longitude"`
	RouteDate                          string  `json:"RouteDate"`
	SRMP                               float64 `json:"SRMP"`
	SRMPAheadBackIndicator             string  `json:"SRMPAheadBackIndicator"`
	StateRouteID                       string  `json:"StateRouteID"`
	StateStructureId                   string  `json:"StateStructureId"`
	VerticalClearanceMaximumFeetInch   string  `json:"VerticalClearanceMaximumFeetInch"`
	VerticalClearanceMaximumInches     int     `json:"VerticalClearanceMaximumInches"`
	VerticalClearanceMinimumFeetInch   string  `json:"VerticalClearanceMinimumFeetInch"`
	VerticalClearanceMinimumInches     int     `json:"VerticalClearanceMinimumInches"`
}

// GetBridgeClearances retrieves all bridge clearance information
func (t *TrafficClient) GetBridgeClearances() ([]BridgeClearance, error) {
	req, err := http.NewRequest(http.MethodGet, getBridgeClearancesAsJsonURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	q := req.URL.Query()
	q.Add(wsdot.ParamTrafficAccessCodeKey, t.wsdot.ApiKey)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.wsdot.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var clearances []BridgeClearance
	if err := json.NewDecoder(resp.Body).Decode(&clearances); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return clearances, nil
}