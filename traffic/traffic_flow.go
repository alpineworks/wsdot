package traffic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"alpineworks.io/wsdot"
)

const (
	getTrafficFlowAsJsonURL = "http://www.wsdot.wa.gov/Traffic/api/TrafficFlow/TrafficFlowREST.svc/GetTrafficFlowAsJson"
)

// TrafficFlowStationLocation represents the location of a traffic flow station
type TrafficFlowStationLocation struct {
	Description string  `json:"Description"`
	Direction   string  `json:"Direction"`
	Latitude    float64 `json:"Latitude"`
	Longitude   float64 `json:"Longitude"`
	MilePost    float64 `json:"MilePost"`
	RoadName    string  `json:"RoadName"`
}

// TrafficFlow represents traffic flow data from a monitoring station
type TrafficFlow struct {
	FlowDataID            int                        `json:"FlowDataID"`
	FlowReadingValue      interface{}                `json:"FlowReadingValue"` // Can be int or string (Unknown, WideOpen, Moderate, Heavy, StopAndGo, NoData)
	FlowStationLocation   TrafficFlowStationLocation `json:"FlowStationLocation"`
	Region                string                     `json:"Region"`
	StationName           string                     `json:"StationName"`
	Time                  string                     `json:"Time"`
}

// GetTrafficFlow retrieves all traffic flow data
func (t *TrafficClient) GetTrafficFlow() ([]TrafficFlow, error) {
	req, err := http.NewRequest(http.MethodGet, getTrafficFlowAsJsonURL, nil)
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

	var trafficFlow []TrafficFlow
	if err := json.NewDecoder(resp.Body).Decode(&trafficFlow); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return trafficFlow, nil
}