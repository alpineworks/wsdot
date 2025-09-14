package traffic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"alpineworks.io/wsdot"
)

const (
	getMountainPassConditionsAsJsonURL = "http://www.wsdot.wa.gov/Traffic/api/MountainPassConditions/MountainPassConditionsREST.svc/GetMountainPassConditionsAsJson"
)

// MountainPassRestriction represents a travel restriction for a mountain pass
type MountainPassRestriction struct {
	RestrictionText   string `json:"RestrictionText"`
	TravelDirection   string `json:"TravelDirection"`
}

// MountainPassCondition represents conditions for a mountain pass
type MountainPassCondition struct {
	DateUpdated                string                  `json:"DateUpdated"`
	ElevationInFeet            int                     `json:"ElevationInFeet"`
	Latitude                   float64                 `json:"Latitude"`
	Longitude                  float64                 `json:"Longitude"`
	MountainPassId             int                     `json:"MountainPassId"`
	MountainPassName           string                  `json:"MountainPassName"`
	RestrictionOne             MountainPassRestriction `json:"RestrictionOne"`
	RestrictionTwo             MountainPassRestriction `json:"RestrictionTwo"`
	RoadCondition              string                  `json:"RoadCondition"`
	TemperatureInFahrenheit    int                     `json:"TemperatureInFahrenheit"`
	TravelAdvisoryActive       bool                    `json:"TravelAdvisoryActive"`
	WeatherCondition           string                  `json:"WeatherCondition"`
}

// GetMountainPassConditions retrieves all mountain pass conditions
func (t *TrafficClient) GetMountainPassConditions() ([]MountainPassCondition, error) {
	req, err := http.NewRequest(http.MethodGet, getMountainPassConditionsAsJsonURL, nil)
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

	var conditions []MountainPassCondition
	if err := json.NewDecoder(resp.Body).Decode(&conditions); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return conditions, nil
}