package traffic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"alpineworks.io/wsdot"
)

const (
	getHighwayAlertsAsJsonURL = "http://www.wsdot.wa.gov/Traffic/api/HighwayAlerts/HighwayAlertsREST.svc/GetAlertsAsJson"
)

// AlertLocation represents a location for highway alerts (reusing RoadwayLocation)
type AlertLocation struct {
	Description string  `json:"Description"`
	Direction   string  `json:"Direction"`
	Latitude    float64 `json:"Latitude"`
	Longitude   float64 `json:"Longitude"`
	MilePost    float64 `json:"MilePost"`
	RoadName    string  `json:"RoadName"`
}

// HighwayAlert represents a highway alert from WSDOT
type HighwayAlert struct {
	AlertID                int           `json:"AlertID"`
	County                 string        `json:"County"`
	EndRoadwayLocation     AlertLocation `json:"EndRoadwayLocation"`
	EndTime                string        `json:"EndTime"`
	EventCategory          string        `json:"EventCategory"`
	EventStatus            string        `json:"EventStatus"`
	ExtendedDescription    string        `json:"ExtendedDescription"`
	HeadlineDescription    string        `json:"HeadlineDescription"`
	LastUpdatedTime        string        `json:"LastUpdatedTime"`
	Priority               string        `json:"Priority"`
	Region                 string        `json:"Region"`
	StartRoadwayLocation   AlertLocation `json:"StartRoadwayLocation"`
	StartTime              string        `json:"StartTime"`
}

// GetHighwayAlerts retrieves all highway alerts
func (t *TrafficClient) GetHighwayAlerts() ([]HighwayAlert, error) {
	req, err := http.NewRequest(http.MethodGet, getHighwayAlertsAsJsonURL, nil)
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

	var alerts []HighwayAlert
	if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return alerts, nil
}