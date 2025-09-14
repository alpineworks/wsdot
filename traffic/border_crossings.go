package traffic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"alpineworks.io/wsdot"
)

const (
	getBorderCrossingsAsJsonURL = "http://www.wsdot.wa.gov/Traffic/api/BorderCrossings/BorderCrossingsREST.svc/GetBorderCrossingsAsJson"
)

// BorderCrossingLocation represents the location details of a border crossing
type BorderCrossingLocation struct {
	Description string  `json:"Description"`
	Direction   string  `json:"Direction"`
	Latitude    float64 `json:"Latitude"`
	Longitude   float64 `json:"Longitude"`
	MilePost    int     `json:"MilePost"`
	RoadName    string  `json:"RoadName"`
}

// BorderCrossing represents a border crossing wait time
type BorderCrossing struct {
	BorderCrossingLocation BorderCrossingLocation `json:"BorderCrossingLocation"`
	CrossingName           string                 `json:"CrossingName"`
	Time                   string                 `json:"Time"`
	WaitTime               int                    `json:"WaitTime"`
}

// GetBorderCrossings retrieves all border crossing wait times
func (t *TrafficClient) GetBorderCrossings() ([]BorderCrossing, error) {
	req, err := http.NewRequest(http.MethodGet, getBorderCrossingsAsJsonURL, nil)
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

	var crossings []BorderCrossing
	if err := json.NewDecoder(resp.Body).Decode(&crossings); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return crossings, nil
}