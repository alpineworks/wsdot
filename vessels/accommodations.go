package vessels

import (
	"encoding/json"
	"fmt"
	"net/http"

	"alpineworks.io/wsdot"
)

const (
	getVesselAccommodationsAsJsonURL = "https://www.wsdot.wa.gov/Ferries/API/Vessels/rest/vesselaccommodations"
)

type VesselAccommodation struct {
	VesselID         int         `json:"VesselID"`
	VesselSubjectID  int         `json:"VesselSubjectID"`
	VesselName       string      `json:"VesselName"`
	VesselAbbrev     string      `json:"VesselAbbrev"`
	Class            VesselClass `json:"Class"`
	CarDeckRestroom  bool        `json:"CarDeckRestroom"`
	CarDeckShelter   bool        `json:"CarDeckShelter"`
	Elevator         bool        `json:"Elevator"`
	ADAAccessible    bool        `json:"ADAAccessible"`
	MainCabinGalley  bool        `json:"MainCabinGalley"`
	MainCabinRestroom bool       `json:"MainCabinRestroom"`
	PublicWifi       bool        `json:"PublicWifi"`
	ADAInfo          string      `json:"ADAInfo"`
	AdditionalInfo   string      `json:"AdditionalInfo"`
}

func (v *VesselsClient) GetVesselAccommodations() ([]VesselAccommodation, error) {
	req, err := http.NewRequest(http.MethodGet, getVesselAccommodationsAsJsonURL, nil)
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

	var accommodations []VesselAccommodation
	if err := json.NewDecoder(resp.Body).Decode(&accommodations); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return accommodations, nil
}

func (v *VesselsClient) GetVesselAccommodationByID(vesselID int) (*VesselAccommodation, error) {
	url := fmt.Sprintf("%s/%d", getVesselAccommodationsAsJsonURL, vesselID)
	
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

	var accommodation VesselAccommodation
	if err := json.NewDecoder(resp.Body).Decode(&accommodation); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &accommodation, nil
}