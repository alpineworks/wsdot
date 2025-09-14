package traffic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"alpineworks.io/wsdot"
)

const (
	getCommercialVehicleRestrictionsAsJsonURL       = "http://www.wsdot.wa.gov/Traffic/api/CVRestrictions/CVRestrictionsREST.svc/GetCommercialVehicleRestrictionsAsJson"
	getCommercialVehicleRestrictionsWithIdAsJsonURL = "http://www.wsdot.wa.gov/Traffic/api/CVRestrictions/CVRestrictionsREST.svc/GetCommercialVehicleRestrictionsWithIdAsJson"
)

// RoadwayLocation represents a location on a roadway
type RoadwayLocation struct {
	Description string  `json:"Description"`
	Direction   string  `json:"Direction"`
	Latitude    float64 `json:"Latitude"`
	Longitude   float64 `json:"Longitude"`
	MilePost    float64 `json:"MilePost"`
	RoadName    string  `json:"RoadName"`
}

// CommercialVehicleRestriction represents a commercial vehicle restriction
type CommercialVehicleRestriction struct {
	BLMaxAxle                          int              `json:"BLMaxAxle"`
	BridgeName                         string           `json:"BridgeName"`
	BridgeNumber                       string           `json:"BridgeNumber"`
	CL8MaxAxle                         int              `json:"CL8MaxAxle"`
	DateEffective                      string           `json:"DateEffective"`
	DateExpires                        string           `json:"DateExpires"`
	DatePosted                         string           `json:"DatePosted"`
	EndRoadwayLocation                 RoadwayLocation  `json:"EndRoadwayLocation"`
	IsDetourAvailable                  bool             `json:"IsDetourAvailable"`
	IsExceptionsAllowed                bool             `json:"IsExceptionsAllowed"`
	IsPermanentRestriction             bool             `json:"IsPermanentRestriction"`
	IsWarning                          bool             `json:"IsWarning"`
	Latitude                           float64          `json:"Latitude"`
	LocationDescription                string           `json:"LocationDescription"`
	LocationName                       string           `json:"LocationName"`
	Longitude                          float64          `json:"Longitude"`
	MaximumGrossVehicleWeightInPounds  int              `json:"MaximumGrossVehicleWeightInPounds"`
	RestrictionComment                 string           `json:"RestrictionComment"`
	RestrictionHeightInInches          int              `json:"RestrictionHeightInInches"`
	RestrictionLengthInInches          int              `json:"RestrictionLengthInInches"`
	RestrictionType                    string           `json:"RestrictionType"`
	RestrictionWeightInPounds          int              `json:"RestrictionWeightInPounds"`
	RestrictionWidthInInches           int              `json:"RestrictionWidthInInches"`
	SAMaxAxle                          int              `json:"SAMaxAxle"`
	StartRoadwayLocation               RoadwayLocation  `json:"StartRoadwayLocation"`
	State                              string           `json:"State"`
	StateRouteID                       string           `json:"StateRouteID"`
	TDMaxAxle                          int              `json:"TDMaxAxle"`
	VehicleType                        string           `json:"VehicleType"`
}

// GetCommercialVehicleRestrictions retrieves all commercial vehicle restrictions
func (t *TrafficClient) GetCommercialVehicleRestrictions() ([]CommercialVehicleRestriction, error) {
	req, err := http.NewRequest(http.MethodGet, getCommercialVehicleRestrictionsAsJsonURL, nil)
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

	var restrictions []CommercialVehicleRestriction
	if err := json.NewDecoder(resp.Body).Decode(&restrictions); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return restrictions, nil
}

// GetCommercialVehicleRestrictionsWithID retrieves all commercial vehicle restrictions with unique IDs
func (t *TrafficClient) GetCommercialVehicleRestrictionsWithID() ([]CommercialVehicleRestriction, error) {
	req, err := http.NewRequest(http.MethodGet, getCommercialVehicleRestrictionsWithIdAsJsonURL, nil)
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

	var restrictions []CommercialVehicleRestriction
	if err := json.NewDecoder(resp.Body).Decode(&restrictions); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return restrictions, nil
}