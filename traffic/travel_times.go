package traffic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"alpineworks.io/wsdot"
)

const (
	getTravelTimesAsJsonURL = "http://www.wsdot.wa.gov/Traffic/api/TravelTimes/TravelTimesREST.svc/GetTravelTimesAsJson"
)

// TravelTimePoint represents a travel time route endpoint
type TravelTimePoint struct {
	Description string  `json:"Description"`
	Direction   string  `json:"Direction"`
	Latitude    float64 `json:"Latitude"`
	Longitude   float64 `json:"Longitude"`
	MilePost    float64 `json:"MilePost"`
	RoadName    string  `json:"RoadName"`
}

// TravelTime represents travel time data for a route segment
type TravelTime struct {
	AverageTime  int             `json:"AverageTime"`
	CurrentTime  int             `json:"CurrentTime"`
	Description  string          `json:"Description"`
	Distance     float64         `json:"Distance"`
	EndPoint     TravelTimePoint `json:"EndPoint"`
	Name         string          `json:"Name"`
	StartPoint   TravelTimePoint `json:"StartPoint"`
	TimeUpdated  string          `json:"TimeUpdated"`
	TravelTimeID int             `json:"TravelTimeID"`
}

// GetTravelTimes retrieves all travel times
func (t *TrafficClient) GetTravelTimes() ([]TravelTime, error) {
	req, err := http.NewRequest(http.MethodGet, getTravelTimesAsJsonURL, nil)
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

	var travelTimes []TravelTime
	if err := json.NewDecoder(resp.Body).Decode(&travelTimes); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return travelTimes, nil
}