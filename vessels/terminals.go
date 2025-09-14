package vessels

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"alpineworks.io/wsdot"
)

const (
	getFareTerminalsURL = "http://www.wsdot.wa.gov/ferries/api/fares/rest/terminals"
)

// FareTerminal represents a terminal in the fares system
type FareTerminal struct {
	TerminalID          int    `json:"TerminalID"`
	TerminalName        string `json:"TerminalName"`
	TerminalAbbrev      string `json:"TerminalAbbrev"`
	Region              int    `json:"Region"`
	SortSeq             int    `json:"SortSeq"`
	Latitude            float64 `json:"Latitude"`
	Longitude           float64 `json:"Longitude"`
	AddressLineOne      string `json:"AddressLineOne"`
	AddressLineTwo      string `json:"AddressLineTwo"`
	City                string `json:"City"`
	State               string `json:"State"`
	ZipCode             string `json:"ZipCode"`
	Country             string `json:"Country"`
	MapLink             string `json:"MapLink"`
	DispGISLat          float64 `json:"DispGISLat"`
	DispGISLong         float64 `json:"DispGISLong"`
	Reservable          bool   `json:"Reservable"`
	CheckTravelAdvisory bool   `json:"CheckTravelAdvisory"`
}

// GetTerminals retrieves all terminals available for fare calculations
func (v *VesselsClient) GetTerminals(tripDate string) ([]FareTerminal, error) {
	url := fmt.Sprintf("%s/%s", getFareTerminalsURL, tripDate)

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

	var terminals []FareTerminal
	if err := json.NewDecoder(resp.Body).Decode(&terminals); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return terminals, nil
}

// GetTerminalMates retrieves terminals that can be reached from a specific terminal
func (v *VesselsClient) GetTerminalMates(tripDate string, terminalID int) ([]FareTerminal, error) {
	url := fmt.Sprintf("http://www.wsdot.wa.gov/ferries/api/fares/rest/terminalmates/%s/%s",
		tripDate, strconv.Itoa(terminalID))

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

	var terminals []FareTerminal
	if err := json.NewDecoder(resp.Body).Decode(&terminals); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return terminals, nil
}