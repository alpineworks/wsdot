package vessels

import (
	"encoding/json"
	"fmt"
	"net/http"

	"alpineworks.io/wsdot"
)

const (
	getVesselStatsAsJsonURL = "https://www.wsdot.wa.gov/Ferries/API/Vessels/rest/vesselstats"
)

type VesselStats struct {
	VesselID                         int         `json:"VesselID"`
	VesselSubjectID                  int         `json:"VesselSubjectID"`
	VesselName                       string      `json:"VesselName"`
	VesselAbbrev                     string      `json:"VesselAbbrev"`
	Class                            VesselClass `json:"Class"`
	VesselNameDesc                   string      `json:"VesselNameDesc"`
	VesselHistory                    string      `json:"VesselHistory"`
	Beam                             string      `json:"Beam"`
	CityBuilt                        string      `json:"CityBuilt"`
	SpeedInKnots                     int         `json:"SpeedInKnots"`
	Draft                            string      `json:"Draft"`
	EngineCount                      int         `json:"EngineCount"`
	Horsepower                       int         `json:"Horsepower"`
	Length                           string      `json:"Length"`
	MaxPassengerCount                int         `json:"MaxPassengerCount"`
	PassengerOnly                    bool        `json:"PassengerOnly"`
	FastFerry                        bool        `json:"FastFerry"`
	PropulsionInfo                   string      `json:"PropulsionInfo"`
	TallDeckClearance                int         `json:"TallDeckClearance"`
	RegDeckSpace                     int         `json:"RegDeckSpace"`
	TallDeckSpace                    int         `json:"TallDeckSpace"`
	Tonnage                          int         `json:"Tonnage"`
	Displacement                     int         `json:"Displacement"`
	YearBuilt                        int         `json:"YearBuilt"`
	YearRebuilt                      int         `json:"YearRebuilt"`
	VesselDrawingImg                 string      `json:"VesselDrawingImg"`
	SolasCertified                   bool        `json:"SolasCertified"`
	PicnicTables                     bool        `json:"PicnicTables"`
	Restroom                         bool        `json:"Restroom"`
	UploadVesselArrivingAnnouncement bool        `json:"UploadVesselArrivingAnnouncement"`
	Galley                           bool        `json:"Galley"`
	VesselWatchUrl                   string      `json:"VesselWatchUrl"`
	Status                           int         `json:"Status"`
	OwnedByWSF                       bool        `json:"OwnedByWSF"`
}

type VesselClass struct {
	ClassID           int    `json:"ClassID"`
	ClassSubjectID    int    `json:"ClassSubjectID"`
	ClassName         string `json:"ClassName"`
	SortSeq           int    `json:"SortSeq"`
	DrawingImg        string `json:"DrawingImg"`
	SilhouetteImg     string `json:"SilhouetteImg"`
	PublicDisplayName string `json:"PublicDisplayName"`
}

func (v *VesselsClient) GetVesselStats() ([]VesselStats, error) {
	req, err := http.NewRequest(http.MethodGet, getVesselStatsAsJsonURL, nil)
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

	var stats []VesselStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return stats, nil
}

func (v *VesselsClient) GetVesselStatsByID(vesselID int) (*VesselStats, error) {
	url := fmt.Sprintf("%s/%d", getVesselStatsAsJsonURL, vesselID)
	
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

	var stats VesselStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &stats, nil
}