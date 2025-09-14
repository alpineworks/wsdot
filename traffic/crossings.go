package traffic

import (
	"alpineworks.io/wsdot"
)

// CrossingsClient provides access to WSDOT Border Crossings API endpoints
type CrossingsClient struct {
	wsdot *wsdot.WSDOTClient
}

// NewCrossingsClient creates a new Border Crossings client
func NewCrossingsClient(wsdotClient *wsdot.WSDOTClient) (*CrossingsClient, error) {
	if wsdotClient == nil {
		return nil, wsdot.ErrNoClient
	}

	return &CrossingsClient{
		wsdot: wsdotClient,
	}, nil
}