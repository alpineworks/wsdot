package traffic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"alpineworks.io/wsdot"
)

const (
	getCurrentWeatherInformationAsJsonURL = "http://www.wsdot.wa.gov/Traffic/api/WeatherInformation/WeatherInformationREST.svc/GetCurrentWeatherInformationAsJson"
)

// WeatherInformation represents current weather information from a station
type WeatherInformation struct {
	BarometricPressure        float64 `json:"BarometricPressure"`
	Latitude                  float64 `json:"Latitude"`
	Longitude                 float64 `json:"Longitude"`
	PrecipitationInInches     float64 `json:"PrecipitationInInches"`
	ReadingTime               string  `json:"ReadingTime"`
	RelativeHumidity          float64 `json:"RelativeHumidity"`
	SkyCoverage               string  `json:"SkyCoverage"`
	StationID                 int     `json:"StationID"`
	StationName               string  `json:"StationName"`
	TemperatureInFahrenheit   float64 `json:"TemperatureInFahrenheit"`
	Visibility                int     `json:"Visibility"`
	WindDirection             float64 `json:"WindDirection"`
	WindDirectionCardinal     string  `json:"WindDirectionCardinal"`
	WindGustSpeedInMPH        float64 `json:"WindGustSpeedInMPH"`
	WindSpeedInMPH            float64 `json:"WindSpeedInMPH"`
}

// GetCurrentWeatherInformation retrieves current weather information from all stations
func (t *TrafficClient) GetCurrentWeatherInformation() ([]WeatherInformation, error) {
	req, err := http.NewRequest(http.MethodGet, getCurrentWeatherInformationAsJsonURL, nil)
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

	var weatherInfo []WeatherInformation
	if err := json.NewDecoder(resp.Body).Decode(&weatherInfo); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return weatherInfo, nil
}