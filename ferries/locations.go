package ferries

import (
	"encoding/json"
	"fmt"
	"net/http"

	"alpineworks.io/wsdot"
)

const (
	getVesselBasicsAsJsonURL    = "https://www.wsdot.wa.gov/Ferries/API/Vessels/rest/vesselbasics"
	getVesselLocationsAsJsonURL = "https://www.wsdot.wa.gov/Ferries/API/Vessels/rest/vessellocations"
)

type VesselBasic struct {
	VesselID        int    `json:"VesselID"`
	VesselSubjectID int    `json:"VesselSubjectID"`
	VesselName      string `json:"VesselName"`
	VesselAbbrev    string `json:"VesselAbbrev"`
	Class           struct {
		ClassID           int    `json:"ClassID"`
		ClassSubjectID    int    `json:"ClassSubjectID"`
		ClassName         string `json:"ClassName"`
		SortSeq           int    `json:"SortSeq"`
		DrawingImg        string `json:"DrawingImg"`
		SilhouetteImg     string `json:"SilhouetteImg"`
		PublicDisplayName string `json:"PublicDisplayName"`
	} `json:"Class"`
	Status     int  `json:"Status"`
	OwnedByWSF bool `json:"OwnedByWSF"`
}

func (f *FerriesClient) GetVesselBasics() ([]VesselBasic, error) {
	req, err := http.NewRequest(http.MethodGet, getVesselBasicsAsJsonURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	q := req.URL.Query()
	q.Add(wsdot.ParamFerriesAccessCodeKey, f.wsdot.ApiKey)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.wsdot.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var vessels []VesselBasic
	if err := json.NewDecoder(resp.Body).Decode(&vessels); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return vessels, nil
}

type VesselLocation struct {
	VesselID                int      `json:"VesselID"`
	VesselName              string   `json:"VesselName"`
	Mmsi                    int      `json:"Mmsi"`
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
	VesselWatchShutID       int      `json:"VesselWatchShutID"`
	VesselWatchShutMsg      string   `json:"VesselWatchShutMsg"`
	VesselWatchShutFlag     string   `json:"VesselWatchShutFlag"`
	VesselWatchStatus       string   `json:"VesselWatchStatus"`
	VesselWatchMsg          string   `json:"VesselWatchMsg"`
}

func (f *FerriesClient) GetVesselLocations() ([]VesselLocation, error) {
	req, err := http.NewRequest(http.MethodGet, getVesselLocationsAsJsonURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	q := req.URL.Query()
	q.Add(wsdot.ParamFerriesAccessCodeKey, f.wsdot.ApiKey)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.wsdot.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var vessels []VesselLocation
	if err := json.NewDecoder(resp.Body).Decode(&vessels); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return vessels, nil
}
