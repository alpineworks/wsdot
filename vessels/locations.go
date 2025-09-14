package vessels

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"alpineworks.io/wsdot"
)

const (
	getTerminalLocationsURL = "http://www.wsdot.wa.gov/ferries/api/terminals/rest/terminallocations"
)

// TerminalLocation represents location information for a terminal
type TerminalLocation struct {
	TerminalID     int     `json:"TerminalID"`
	TerminalName   string  `json:"TerminalName"`
	TerminalAbbrev string  `json:"TerminalAbbrev"`
	Latitude       float64 `json:"Latitude"`
	Longitude      float64 `json:"Longitude"`
	AddressLineOne string  `json:"AddressLineOne"`
	AddressLineTwo string  `json:"AddressLineTwo"`
	City           string  `json:"City"`
	State          string  `json:"State"`
	ZipCode        string  `json:"ZipCode"`
	Country        string  `json:"Country"`
	MapLink        string  `json:"MapLink"`
	DispGISLat     float64 `json:"DispGISLat"`
	DispGISLong    float64 `json:"DispGISLong"`
}

// GetTerminalLocations retrieves location information for all terminals
func (v *VesselsClient) GetTerminalLocations() ([]TerminalLocation, error) {
	req, err := http.NewRequest(http.MethodGet, getTerminalLocationsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	q := req.URL.Query()
	q.Add(wsdot.ParamFerriesAccessCodeKey, v.wsdot.ApiKey)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.wsdot.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var locations []TerminalLocation
	if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return locations, nil
}

// GetTerminalLocationByID retrieves location information for a specific terminal
func (v *VesselsClient) GetTerminalLocationByID(terminalID int) (*TerminalLocation, error) {
	url := fmt.Sprintf("%s/%s", getTerminalLocationsURL, strconv.Itoa(terminalID))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	q := req.URL.Query()
	q.Add(wsdot.ParamFerriesAccessCodeKey, v.wsdot.ApiKey)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.wsdot.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var location TerminalLocation
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &location, nil
}