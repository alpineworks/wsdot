package traffic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"alpineworks.io/wsdot"
)

const (
	getTollRatesAsJsonURL = "http://www.wsdot.wa.gov/Traffic/api/TollRates/TollRatesREST.svc/GetTollRatesAsJson"
)

// TollRate represents toll rate information for a route segment
type TollRate struct {
	CurrentMessage      string  `json:"CurrentMessage"`
	CurrentToll         float64 `json:"CurrentToll"`
	EndLatitude         float64 `json:"EndLatitude"`
	EndLocationName     string  `json:"EndLocationName"`
	EndLongitude        float64 `json:"EndLongitude"`
	EndMilepost         float64 `json:"EndMilepost"`
	StartLatitude       float64 `json:"StartLatitude"`
	StartLocationName   string  `json:"StartLocationName"`
	StartLongitude      float64 `json:"StartLongitude"`
	StartMilepost       float64 `json:"StartMilepost"`
	StateRoute          string  `json:"StateRoute"`
	TimeUpdated         string  `json:"TimeUpdated"`
	TravelDirection     string  `json:"TravelDirection"`
	TripName            string  `json:"TripName"`
}

// GetTollRates retrieves all toll rates
func (t *TrafficClient) GetTollRates() ([]TollRate, error) {
	req, err := http.NewRequest(http.MethodGet, getTollRatesAsJsonURL, nil)
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

	var tollRates []TollRate
	if err := json.NewDecoder(resp.Body).Decode(&tollRates); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return tollRates, nil
}