package vessels

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"alpineworks.io/wsdot"
)

const (
	getVesselHistoryAsJsonURL = "https://www.wsdot.wa.gov/Ferries/API/Vessels/rest/vesselhistory"
)

type VesselHistory struct {
	VesselID                int      `json:"VesselID"`
	VesselName              string   `json:"VesselName"`
	DepartingTerminalID     int      `json:"DepartingTerminalID"`
	DepartingTerminalName   string   `json:"DepartingTerminalName"`
	DepartingTerminalAbbrev string   `json:"DepartingTerminalAbbrev"`
	ArrivingTerminalID      int      `json:"ArrivingTerminalID"`
	ArrivingTerminalName    string   `json:"ArrivingTerminalName"`
	ArrivingTerminalAbbrev  string   `json:"ArrivingTerminalAbbrev"`
	Latitude                float64  `json:"Latitude"`
	Longitude               float64  `json:"Longitude"`
	Speed                   float64  `json:"Speed"`
	Heading                 int      `json:"Heading"`
	InService               bool     `json:"InService"`
	AtDock                  bool     `json:"AtDock"`
	LeftDock                string   `json:"LeftDock"`
	Eta                     string   `json:"Eta"`
	EtaBasis                string   `json:"EtaBasis"`
	ScheduledDeparture      string   `json:"ScheduledDeparture"`
	OpRouteAbbrev           []string `json:"OpRouteAbbrev"`
	VesselPositionNum       int      `json:"VesselPositionNum"`
	SortSeq                 int      `json:"SortSeq"`
	ManagedBy               int      `json:"ManagedBy"`
	TimeStamp               string   `json:"TimeStamp"`
}

func (v *VesselsClient) GetVesselHistory() ([]VesselHistory, error) {
	req, err := http.NewRequest(http.MethodGet, getVesselHistoryAsJsonURL, nil)
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

	var history []VesselHistory
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return history, nil
}

func (v *VesselsClient) GetVesselHistoryByNameAndDateRange(vesselName string, dateStart, dateEnd time.Time) ([]VesselHistory, error) {
	// Format dates as YYYY-MM-DD
	startStr := dateStart.Format("2006-01-02")
	endStr := dateEnd.Format("2006-01-02")
	
	url := fmt.Sprintf("%s/%s/%s/%s", getVesselHistoryAsJsonURL, vesselName, startStr, endStr)
	
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

	var history []VesselHistory
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return history, nil
}