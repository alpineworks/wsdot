package vessels

import (
	"alpineworks.io/wsdot"
)

// FaresClient provides access to WSDOT Fares API endpoints
type FaresClient struct {
	wsdot *wsdot.WSDOTClient
}

// NewFaresClient creates a new Fares client
func NewFaresClient(wsdotClient *wsdot.WSDOTClient) (*FaresClient, error) {
	if wsdotClient == nil {
		return nil, wsdot.ErrNoClient
	}

	return &FaresClient{
		wsdot: wsdotClient,
	}, nil
}