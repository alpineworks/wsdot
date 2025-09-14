package traffic

import (
	"alpineworks.io/wsdot"
)

// TrafficClient provides access to WSDOT Traffic API endpoints
type TrafficClient struct {
	wsdot *wsdot.WSDOTClient
}

// NewTrafficClient creates a new Traffic client
func NewTrafficClient(wsdotClient *wsdot.WSDOTClient) (*TrafficClient, error) {
	if wsdotClient == nil {
		return nil, wsdot.ErrNoClient
	}

	return &TrafficClient{
		wsdot: wsdotClient,
	}, nil
}

// GetWSDOTClient returns the underlying WSDOT client
func (t *TrafficClient) GetWSDOTClient() *wsdot.WSDOTClient {
	return t.wsdot
}