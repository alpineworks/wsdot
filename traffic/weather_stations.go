package traffic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"alpineworks.io/wsdot"
)

const (
	getWeatherStationsAsJsonURL = "http://www.wsdot.wa.gov/Traffic/api/WeatherStations/WeatherStationsREST.svc/GetWeatherStationsAsJson"
)

// WeatherStation represents a weather monitoring station
type WeatherStation struct {
	Description string  `json:"Description"`
	Latitude    float64 `json:"Latitude"`
	Longitude   float64 `json:"Longitude"`
	StationID   int     `json:"StationID"`
	StationName string  `json:"StationName"`
}

// GetWeatherStations retrieves all weather station information
func (t *TrafficClient) GetWeatherStations() ([]WeatherStation, error) {
	req, err := http.NewRequest(http.MethodGet, getWeatherStationsAsJsonURL, nil)
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

	var stations []WeatherStation
	if err := json.NewDecoder(resp.Body).Decode(&stations); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return stations, nil
}