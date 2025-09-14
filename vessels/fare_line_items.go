package vessels

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"alpineworks.io/wsdot"
)

const (
	getFareLineItemsURL      = "http://www.wsdot.wa.gov/ferries/api/fares/rest/farelineitems"
	getFareLineItemsBasicURL = "http://www.wsdot.wa.gov/ferries/api/fares/rest/farelineitemsbasic"
)

// FareLineItem represents a fare line item for a specific route and trip
type FareLineItem struct {
	FareLineItemID          int     `json:"FareLineItemID"`
	FareTypeID              int     `json:"FareTypeID"`
	FareTypeDescription     string  `json:"FareTypeDescription"`
	FareLineItemDescription string  `json:"FareLineItemDescription"`
	Cost                    float64 `json:"Cost"`
	SortOrder               int     `json:"SortOrder"`
	MaxVehicleLength        int     `json:"MaxVehicleLength"`
	MaxVehicleHeight        int     `json:"MaxVehicleHeight"`
	VehicleClass            string  `json:"VehicleClass"`
	IsVehicle               bool    `json:"IsVehicle"`
	RestrictionText         string  `json:"RestrictionText"`
	Discountable            bool    `json:"Discountable"`
}

// BasicFareLineItem represents a simplified fare line item
type BasicFareLineItem struct {
	FareLineItemID          int     `json:"FareLineItemID"`
	FareLineItemDescription string  `json:"FareLineItemDescription"`
	Cost                    float64 `json:"Cost"`
	SortOrder               int     `json:"SortOrder"`
	IsVehicle               bool    `json:"IsVehicle"`
}

// GetFareLineItems retrieves detailed fare line items for a specific route and trip
func (v *VesselsClient) GetFareLineItems(tripDate, departingTerminalID, arrivingTerminalID string, roundTrip bool) ([]FareLineItem, error) {
	roundTripStr := "false"
	if roundTrip {
		roundTripStr = "true"
	}

	url := fmt.Sprintf("%s/%s/%s/%s/%s", getFareLineItemsURL, tripDate, departingTerminalID, arrivingTerminalID, roundTripStr)

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

	var fareItems []FareLineItem
	if err := json.NewDecoder(resp.Body).Decode(&fareItems); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return fareItems, nil
}

// GetFareLineItemsBasic retrieves basic fare line items for a specific route and trip
func (v *VesselsClient) GetFareLineItemsBasic(tripDate, departingTerminalID, arrivingTerminalID string, roundTrip bool) ([]BasicFareLineItem, error) {
	roundTripStr := "false"
	if roundTrip {
		roundTripStr = "true"
	}

	url := fmt.Sprintf("%s/%s/%s/%s/%s", getFareLineItemsBasicURL, tripDate, departingTerminalID, arrivingTerminalID, roundTripStr)

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

	var fareItems []BasicFareLineItem
	if err := json.NewDecoder(resp.Body).Decode(&fareItems); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return fareItems, nil
}

// GetFareTotal calculates the total fare for specified items and quantities
func (v *VesselsClient) GetFareTotal(tripDate, departingTerminalID, arrivingTerminalID string, roundTrip bool, fareLineItemID int, quantity int) (*FareTotalResult, error) {
	roundTripStr := "false"
	if roundTrip {
		roundTripStr = "true"
	}

	url := fmt.Sprintf("http://www.wsdot.wa.gov/ferries/api/fares/rest/faretotals/%s/%s/%s/%s/%s/%s",
		tripDate, departingTerminalID, arrivingTerminalID, roundTripStr,
		strconv.Itoa(fareLineItemID), strconv.Itoa(quantity))

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

	var result FareTotalResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

// FareTotalResult represents the result of a fare total calculation
type FareTotalResult struct {
	TotalCost           float64 `json:"TotalCost"`
	Currency            string  `json:"Currency"`
	TaxAmount           float64 `json:"TaxAmount"`
	FareLineItemID      int     `json:"FareLineItemID"`
	Quantity            int     `json:"Quantity"`
	UnitCost            float64 `json:"UnitCost"`
	Description         string  `json:"Description"`
}